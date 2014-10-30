// this is a poor quality bazil implementation of a loopback filesystem
// +build bazil

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"fmt"
	"io/ioutil"
	"os"
)

const (
	loopfs_rdev = 0x70F5 //LOop-FS
)

//ok
type looperfs struct {
	d *dump
}

type looperdir struct {
	d *dump
	name string
}

type looperfile struct {
	d *dump
	name string
	f    *os.File
}

func (l looperfs) Root() (fs.Node, fuse.Error) {
	return looperdir{name: ".", d: l.d}, nil
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
		return looperdir{name: l.name + "/" + name, d: l.d}, nil
	} else {
		return looperfile{name: l.name + "/" + name, d: l.d}, nil
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

func (l looperfile) Open(req *fuse.OpenRequest, resp *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	file, err := ioutil.ReadFile(l.name)

	return fs.DataHandle(file), err

}

func (l looperfile) Write(req *fuse.WriteRequest, resp *fuse.WriteResponse, intr fs.Intr) fuse.Error {

	size, err := l.f.Write(req.Data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	resp.Size = size

	//	fmt.Println("WRITE at ", req.Offset)
	return nil
}

func (l looperdir) Mkdir(req *fuse.MkdirRequest, intr fs.Intr) (fs.Node, fuse.Error) {
	os.Mkdir(req.Name, req.Mode)
	return looperdir{name: l.name + "/" + req.Name}, nil
}

func (l looperdir) Remove(req *fuse.RemoveRequest, intr fs.Intr) fuse.Error {
	os.Remove(req.Name)
	return nil
}

func (l looperfile) Remove(req *fuse.RemoveRequest, intr fs.Intr) fuse.Error {
	os.Remove(req.Name)
	return nil
}

func (l looperfile) Release(req *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
	fmt.Println("release at ")
	return nil
}

func (l looperdir) Create(req *fuse.CreateRequest, resp *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
	fname := l.name + "/" + req.Name

	fi, err := os.Create(fname)
	if err != nil {
		return nil, nil, err
	}
	//	l.f = fi
	f := looperfile{name: fname, f: fi}

	return f, f, nil
}

//TODO: file: setattr,lock/unlock
