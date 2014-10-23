// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func (f *Ffs) monut() (e error) {
	pathFs := pathfs.NewPathNodeFs(pathfs.NewLoopbackFileSystem("foo"), nil)
	connector := nodefs.NewFileSystemConnector(pathFs.Root(), nil)
	f.be.gc, e = fuse.NewServer(fuse.NewRawFileSystem(connector.RawFS()),
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

type fbackend struct {
	gc *fuse.Server
}
