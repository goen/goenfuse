// this is a poor quality bazil implementation of a bin filesystem
// +build bazil

package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"os"
	"time"
)

const (
	inodeoffset = 5
)

func tapcontext() fs.FS {
	return tapperfs{r: tapperrootnode{}}
}

//ok
func generic_attr() fuse.Attr {
	u := time.Unix(0, 0)
	return fuse.Attr{
		Atime: u, Mtime: u, Ctime: u, Crtime: u,
		Uid:   uint32(os.Geteuid()),
		Gid:   uint32(os.Getegid()),
		Rdev:  0xB1F5,     //BIn-FS
		Flags: 0x00121012, //don't modify
	}
}

// tapperFS dirs: root dir & the various dirs
//

//ok
type tapperfs struct {
	r tapperrootnode
}

//ok
type tapperrootnode struct {}

//ok
type tapperlooppipe struct {}
type tappertracepipe struct {}

//ok

// get fs root node
func (s tapperfs) Root() (fs.Node, fuse.Error) {
	return s.r, nil
}

func (tapperrootnode) Attr() fuse.Attr {
	a := generic_attr()
	a.Inode = 2
	a.Size = 4096
	a.Blocks = 8
	a.Mode = os.ModeDir | 0555
	a.Nlink = 8 // correct?//FIXME
	return a
}

//ok
func (s tapperrootnode) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {


	switch (name) {
	case "loop":
		return tapperlooppipe{}, nil
	case "trace":
		return tappertracepipe{}, nil

	default:
	return nil, fuse.ENOENT
	}
}

//ok
func (s tapperrootnode) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var dirz [2]fuse.Dirent
	foffset := uint64(inodeoffset)

	//DT_FIFO uppercase man
	dirz[0] = fuse.Dirent{Inode: foffset, Name: "loop", Type: fuse.DT_FIFO}
	dirz[1] = fuse.Dirent{Inode: foffset+1, Name: "trace", Type: fuse.DT_FIFO}

	sdirs := dirz[0:]

	return sdirs, nil
}

func (tapperlooppipe) Attr() fuse.Attr {
	return fuse.Attr{Mode: os.ModeNamedPipe | 0777}
}

func (tappertracepipe) Attr() fuse.Attr {
	return fuse.Attr{Mode: os.ModeNamedPipe | 0777}
}
