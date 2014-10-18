package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"io/ioutil"
	"os"
)

const (
	loopfs_rdev = 0x70F5 //LOop-FS
)

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

func (l looperdir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if l.name == "." && (name == mpoint_gloop || name == mpoint_gbin) {
		return nil, fuse.ENOENT
	}

	fi, err := os.Lstat(l.name + "/" + name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	if fi.IsDir() {
		return looperdir{name: l.name + "/" + name}, nil
	} else {
		return looperfile{name: l.name + "/" + name}, nil
	}

	return nil, fuse.ENOENT
}

func (l looperdir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	fi, err := ioutil.ReadDir(l.name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	var dirz []fuse.Dirent

	//FIXME: add files: . ..

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

func (l looperdir) Attr() fuse.Attr {
	b := looperfile{name: l.name}.Attr()

	b.Mode = os.ModeDir | b.Mode
	return b
}

func (l looperfile) Attr() fuse.Attr {
	fi, err := os.Lstat(l.name)
	if err != nil {
		return fuse.Attr{}
	}

	s := StatSys(fi)

	b := fuse.Attr{
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

	b.Rdev = loopfs_rdev
	b.Inode ^= 0x7fff

	return b
}
