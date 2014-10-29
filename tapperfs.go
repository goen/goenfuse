// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

const (
	inodeoffset = 5
)

func tapcontext() nodefs.Node {
	return tapper_root{}
}

type tapperlooppipe struct {nodefs.Node}
type tappertracepipe struct {nodefs.Node}

type tapper_root struct {
	nodefs.Node
}

func newtapperlooppipe() *tapperlooppipe {
	return &tapperlooppipe{Node: nodefs.NewDefaultNode()}
}

func newtappertracepipe() *tappertracepipe {
	return &tappertracepipe{Node: nodefs.NewDefaultNode()}
}

func (r tapper_root) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	var dirz [2]fuse.DirEntry

	dirz[0] = fuse.DirEntry{Name: "loop", Mode: 0777}
	dirz[1] = fuse.DirEntry{Name: "trace", Mode: 0777}

	sdirs := dirz[0:]

	return sdirs, fuse.OK
}

func (r tapper_root) Lookup(out *fuse.Attr, name string, context *fuse.Context) (node *nodefs.Inode, code fuse.Status) {
	S_IFIFO := uint32(0x1000)

	switch (name) {
	case "loop":
		out.Mode = S_IFIFO | 0777
		ch := r.Inode().NewChild(name, false, newtapperlooppipe())
		return ch, fuse.OK
	case "trace":
		out.Mode = S_IFIFO | 0777
		ch := r.Inode().NewChild(name, false, newtappertracepipe())
		return ch, fuse.OK

	default:
	return nil, fuse.ENOENT
	}
}

func (tapperlooppipe) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {
	out.Mode = 0x1000 | 0777

	return fuse.OK
}

func (tappertracepipe) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {
	out.Mode = 0x1000 | 0777

	return fuse.OK
}

func (tapper_root) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFDIR | 0755

	return fuse.OK
}
