// +build bazil

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func (f *Ffs) monut() (e error) {
	f.be.bc, e = fuse.Mount(f.dir)
	return e
}

func (f *Ffs) unmount() (err error) {
	return fuse.Unmount(f.dir)
}

func (f Ffs) check() error {
	<-f.be.bc.Ready
	if err := f.be.bc.MountError; err != nil {
		return err
	}
	return nil
}

func (f Ffs) serve() {
	fs.Serve(f.be.bc, f.be.bs)
}

func destory(f Ffs) {
	f.be.bc.Close()
}

type fbackend struct {
	bs fs.FS
	bc *fuse.Conn
}
