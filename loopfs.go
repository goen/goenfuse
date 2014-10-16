package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	_ "fmt"
	"io/ioutil"
	"os"
	"time"
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

func (looperdir) Attr() fuse.Attr {
	a := lgeneric_attr()
	a.Inode = 1
	a.Size = 4096
	a.Blocks = 8
	a.Mode = os.ModeDir | 0555
	a.Nlink = 8 // correct?//FIXME
	return a
}

func (l looperdir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {

	fi, err := ioutil.ReadDir(l.name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	for i := range fi {
		if name == fi[i].Name() {
			return looperfile{}, nil
		}
	}

	return nil, fuse.ENOENT
}

func (l looperdir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {

	fi, err := ioutil.ReadDir(l.name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	var dirz []fuse.Dirent

	for i := range fi {
		node := fuse.Dirent{Inode: 2, Name: fi[i].Name(), Type: fuse.DT_File}
		dirz = append(dirz, node)
	}

	return dirz, nil
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
