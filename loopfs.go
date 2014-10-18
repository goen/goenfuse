package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	loopfs_rdev = 0x70F5 //LOop-FS
)

//ok
func lgeneric_attr() fuse.Attr {
	u := time.Now()
	return fuse.Attr{
		Atime: u, Mtime: u, Ctime: u, Crtime: u,
		Uid:   uint32(os.Geteuid()),
		Gid:   uint32(os.Getegid()),
		Rdev:  0x70F5, //LOop-FS
		Flags: 0,
	}
}

func (looperdir) Attr() fuse.Attr {
	a := lgeneric_attr()
	a.Inode = 1
	a.Size = 4096
	a.Blocks = 8
	a.Mode = os.ModeDir | 0555
	a.Nlink = 8 // correct?//FIXME
	return a
}

func (looperfile) Attr() fuse.Attr {
	a := lgeneric_attr()
	a.Inode = 2
	a.Size = 4096
	a.Blocks = 8
	a.Mode = 0555
	a.Nlink = 1 // correct?//FIXME
	return a
}

func to(s Stat) fuse.Attr {
	return fuse.Attr{
		Inode:  s.Inode(),
		Size:   s.Size_(),
		Blocks: s.Blocks_(),
		Atime:  s.Atime(),
		Mtime:  s.Mtime(),
		Ctime:  s.Ctime(),
		Crtime: s.Crtime(),
		Mode:   s.Mode_(),
		Nlink:  s.Nlink_(),
		Uid:    s.Uid,
		Gid:    s.Gid,
		Rdev:   s.Rdev_(),
		Flags:  s.Flags(),
	}
}

func chinode(f fuse.Attr) fuse.Attr {
	f.Inode ^= 0x7FFF
	return f
}

func putrdev(f fuse.Attr) fuse.Attr {
	f.Rdev = loopfs_rdev
	return f
}

//ok
type looperfs struct {
}

type looperdir struct {
	name string
}

type looperfile struct {
	name string
}

// get fs root node
func (looperfs) Root() (fs.Node, fuse.Error) {
	return looperdir{name: "."}, nil
}

/*
func (l looperdir) Attr() fuse.Attr {
	return looperfile{name: l.name}.Attr()
}
*/
func (l looperdir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if l.name == "." && (name == mpoint_gloop || name == mpoint_gbin) {
		return nil, fuse.ENOENT
	}

	fi, err := os.Lstat(l.name + "/" + name)
	if err != nil {
		return looperdir{name: l.name}, nil
		//		return nil, fuse.ENOENT
	}

	if fi.IsDir() {
		return looperdir{name: l.name + "/" + name}, nil
	} else {
		return looperfile{name: l.name + "/" + name}, nil
	}

	//	return looperfile{name: l.name + "/" + name}, nil
	return nil, fuse.ENOENT
}

func (l looperdir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {

	fmt.Println("reading dir:", l.name)

	fi, err := ioutil.ReadDir(l.name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	var dirz []fuse.Dirent

	for i := range fi { //FIXME: inodez, typez
		name := fi[i].Name()

		if name == mpoint_gloop || name == mpoint_gbin {
			continue
		}

		node := fuse.Dirent{Inode: 2, Name: name, Type: fuse.DT_File}
		dirz = append(dirz, node)
	}

	return dirz, nil
}

/*
func (l looperfile) Attr() fuse.Attr {
	a, err := os.Lstat(l.name)
	if err != nil {
		return fuse.Attr{}
	}

	return putrdev(chinode(to(StatSys(a))))
}
*/
