// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

//

func tapcontext(i interface{}, z interface{}) interface{} {
	return nil
}

func loopcontext() interface{} {
	return nil
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

func (f *Ffs) monut() (e error) {
	return nil
}

func (f *Ffs) putcontext() (e error) {
	//	switch

	pathFs := pathfs.NewPathNodeFs(pathfs.NewLoopbackFileSystem("foo"+f.dir), nil)
	con := nodefs.NewFileSystemConnector(pathFs.Root(), nil)
	f.be.gc, e = fuse.NewServer(fuse.NewRawFileSystem(con.RawFS()),
		f.dir, &fuse.MountOptions{SingleThreaded: true})
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

type stuffer interface{}

type fbackend struct {
	gc *fuse.Server
}
