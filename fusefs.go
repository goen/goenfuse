// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	//	"github.com/hanwen/go-fuse/fuse/pathfs"
)

//

func tapcontext(i interface{}, z interface{}) interface{} {
	return int(1)
}

func loopcontext() interface{} {
	return int(0)
}

////////////////////////////////////////////

type foo struct {
	nodefs.Node
}
type bar struct {
	nodefs.Node
}

// this is not used in

type looperfs struct {
	path string
}
type tapperrootnode struct {
	itemz interface{}
	s     interface{}
}
type tapperfs struct{ r tapperrootnode }

////////////////////////////////////////////

func (f *Ffs) monut() (e error) {
	return nil
}

func (f *Ffs) putcontext() (e error) {
	what := f.stuff.(int)

	var my nodefs.Node

	if what == 0 {

		my = &foo{nodefs.NewDefaultNode()}
		/*
			pathFs := pathfs.NewPathNodeFs(pathfs.NewLoopbackFileSystem("foo"+f.dir), nil)
		*/

	} else {

		my = &bar{nodefs.NewDefaultNode()}
		/*
			pathFs := pathfs.NewPathNodeFs(pathfs.NewLoopbackFileSystem("foo"+f.dir), nil)

		*/

	}

	con := nodefs.NewFileSystemConnector(my, nil)
	raw := fuse.NewRawFileSystem(con.RawFS())
	optz := &fuse.MountOptions{SingleThreaded: true}

	f.be.gc, e = fuse.NewServer(raw, f.dir, optz)
	//	f.be.gc, _, e = nodefs.MountRoot(f.dir, /*root node*/, nil)

	//(*fuse.Server, *FileSystemConnector, error)

	return e
}

func (f *Ffs) unmount() (err error) {
	f.be.gc.Unmount()
	return nil
}

func (f Ffs) serve() {
	f.be.gc.Serve()
}

func (f Ffs) check() error {
	f.be.gc.WaitMount()
	return nil
}

func destory(f Ffs) {
	//XXX go-fuse destructor
}

const (
	bazilfs = false
)

type stuffer interface {
}

type fbackend struct {
	gc *fuse.Server
}
