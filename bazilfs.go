// this is the bazil glue
// +build bazil

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func (Ffs) umt3() int {
	//unmount retries
	return 100
}

func (f *Ffs) monut() (e error) {
	f.be.c, e = fuse.Mount(f.dir)

	return e
}

func (f *Ffs) putcontext() (e error) {
	if f.d == nil {
		f.be.s = fs.FS(tapperfs{r: tapperrootnode{}})
	} else {
		f.be.s = fs.FS(looperfs{d: f.d})
	}
	return nil
}

func (f *Ffs) unmount() (err error) {
	return fuse.Unmount(f.dir)
}

func (f Ffs) check() error {
	<-f.be.c.Ready
	if err := f.be.c.MountError; err != nil {
		return err
	}
	return nil
}

func (f Ffs) serve() {
	fs.Serve(f.be.c, f.be.s)
}

func destory(f Ffs) {
	f.be.c.Close()
}

type fbackend struct {
	s fs.FS
	c *fuse.Conn
}
