// +build bazil

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func (f *Ffs) monut() (e error) {
	f.be.c, e = fuse.Mount(f.dir)
	f.be.s = f.stuff

	return e
}

func (f *Ffs) mount() (e error) {
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

const (
	bazilfs = true
)

type fbackend struct {
	s fs.FS
	c *fuse.Conn
}
