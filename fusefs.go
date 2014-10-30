// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

////////////////////////////////////////////

func (Ffs) umt3() int {
	//unmount retries
	return 3
}

func (f *Ffs) monut() (e error) {
	return nil
}

func (f *Ffs) putcontext() (e error) {
	var my nodefs.Node
	var optz fuse.MountOptions

	if f.d == nil {
		my = &tapper_root{Node: nodefs.NewDefaultNode()}

		optz.SingleThreaded = true
	} else {

		finalFs := NewLooperFileSystem(".", f.d)
		pathFs := pathfs.NewPathNodeFs(finalFs, nil)

		my = pathFs.Root()
		optz.SingleThreaded = false
	}

	con := nodefs.NewFileSystemConnector(my, nil)
	raw := fuse.NewRawFileSystem(con.RawFS())

	f.be.gc, e = fuse.NewServer(raw, f.dir, &optz)

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
