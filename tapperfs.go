// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"fmt"
	"os"
)

const (
	foffset = 2
)

func tapcontext(i [][]string, z *self, pathz *[]string) nodefs.Node {
	return tapper_root{itemz: i, self: z, pathz: pathz}
}

type tapper_real_pathsnode struct {
	nodefs.Node
	pathz *[]string
}

type tapper_root struct {
	nodefs.Node
	pathz *[]string
	itemz [][]string // len(itemz) = 1 + maximum name
	*self
}

type tappertrackernode struct {
	nodefs.Node
	*self
	f *os.File
}

type tapperdirnode struct {
	nodefs.Node
	i     uint64 //name = i, inode = i + inodeoffset
	itemz []string
}

type tapperbinlink struct {
	nodefs.Node
	inode uint64
}

func nrn(s *self) *tappertrackernode {
	name, _ := s.get()
	f, err := os.Open(name)
	if err != nil {
		f = nil
	}

	return &tappertrackernode{Node: nodefs.NewDefaultNode(), self: s, f: f}
}

func gfd(pathz *[]string) *tapper_real_pathsnode {
	return &tapper_real_pathsnode{Node: nodefs.NewDefaultNode(), pathz: pathz}
}

func ndn(i int, itemz []string) *tapperdirnode {
	return &tapperdirnode{Node: nodefs.NewDefaultNode(), i: uint64(i), itemz: itemz}
}

func nln(inode uint64) *tapperbinlink {
	return &tapperbinlink{Node: nodefs.NewDefaultNode(), inode: uint64(inode)}
}

func (s tapperdirnode) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	var foobar []fuse.DirEntry

	for i := range s.itemz {
		item := fuse.DirEntry{Name: s.itemz[i], Mode: 0555}
		foobar = append(foobar, item)
	}
	return foobar, fuse.OK
}

func (s tapperdirnode) Lookup(out *fuse.Attr, name string, context *fuse.Context) (node *nodefs.Inode, code fuse.Status) {
	ibase := (s.i << 18) + 128

	// TODO: binary search
	for i := range s.itemz {
		if name == s.itemz[i] {

			ch := s.Inode().NewChild(name, false, nln(ibase+uint64(i)))

			return ch, fuse.OK
		}
	}

	return nil, fuse.ENOSYS
}

func (r tapper_root) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	var dirz [100 + foffset]fuse.DirEntry

	dirz[0] = fuse.DirEntry{Name: "tracker", Mode: 0555}
	dirz[1] = fuse.DirEntry{Name: "abspaths", Mode: 0555}

	end := int(len(r.itemz))
	if end >= 100 {
		end = 100
	}

	for i := 0; i < end; i++ {
		dirz[i+foffset].Mode = uint32(os.ModeDir | 0555)
		dirz[i+foffset].Name = fmt.Sprintf("%02d", i)
	}

	sdirs := dirz[0 : end+foffset]

	return sdirs, fuse.OK
}

func (r tapper_root) Lookup(out *fuse.Attr, name string, context *fuse.Context) (node *nodefs.Inode, code fuse.Status) {
	if name == "tracker" {
		out.Mode = fuse.S_IFLNK
		ch := r.Inode().NewChild(name, false, nrn(r.self))
		return ch, fuse.OK
	}

	if name == "abspaths" {
		out.Mode = fuse.S_IFREG
		ch := r.Inode().NewChild(name, false, gfd(r.pathz))
		return ch, fuse.OK
	}

	var i int

	n, err := fmt.Sscanf(name, "%02d", &i)

	if (err != nil) || (n != 1) || (i >= len(r.itemz)) {
		return nil, fuse.ENOENT
	}

	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	ch := r.Inode().NewChild(name, true, ndn(i, r.itemz[i]))

	return ch, fuse.OK
}
func (tapperdirnode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFDIR | 0755

	return fuse.OK
}
func (t tappertrackernode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {
	_, size := t.self.get()
	out.Mode = fuse.S_IFREG | 0555
	out.Size = uint64(size)
	return fuse.OK
}
func (t tappertrackernode) Open(flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	return nodefs.NewDefaultFile(), fuse.OK
}
func (t tappertrackernode) Read(file nodefs.File, dest []byte, off int64, context *fuse.Context) (fuse.ReadResult, fuse.Status) {
	if t.f != nil {
		t.f.ReadAt(dest, off)
		return fuse.ReadResultData(dest), fuse.OK
	}
	return nil, fuse.ENOSYS
}

func (tapperbinlink) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFLNK | 0555
	out.Size = 10
	return fuse.OK
}

func (tapperbinlink) Readlink(c *fuse.Context) ([]byte, fuse.Status) {
	return []byte("../tracker"), fuse.OK
}

func (tapper_root) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFDIR | 0755

	return fuse.OK
}
