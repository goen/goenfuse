package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"fmt"
	"os"
	"time"
)

const (
	inodeoffset = 5
)

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
type tapperrootnode struct {
	itemz [][]string // len(itemz) = 1 + maximum name
}

//ok
type tappertrackernode struct {
}

//ok
type tapperdirnode struct {
	i     uint64 //name = i, inode = i + inodeoffset
	itemz []string
}

type tapperbinlink struct {
	inode uint64
}

//ok

// get fs root node
func (s tapperfs) Root() (fs.Node, fuse.Error) {
	return s.r, nil
}

func (tapperrootnode) Attr() fuse.Attr {
	a := generic_attr()
	a.Inode = 1
	a.Size = 4096
	a.Blocks = 8
	a.Mode = os.ModeDir | 0555
	a.Nlink = 8 // correct?//FIXME
	return a
}

//ok
func (s tapperrootnode) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if name == "tracker" || name == ".track" || name == ".untrack" {
		return tappertrackernode{}, nil
	}

	var i int

	n, err := fmt.Sscanf(name, "%02d", &i)

	if (err != nil) || (n != 1) || (i >= len(s.itemz)) {
		return nil, fuse.ENOENT
	}

	return tapperdirnode{i: uint64(i), itemz: s.itemz[i]}, nil
}

//ok
func (s tapperrootnode) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var dirz [103]fuse.Dirent
	foffset := 3

	dirz[0] = fuse.Dirent{Inode: 2, Name: "tracker", Type: fuse.DT_File}
	dirz[1] = fuse.Dirent{Inode: 3, Name: ".track", Type: fuse.DT_File}
	dirz[2] = fuse.Dirent{Inode: 4, Name: ".untrack", Type: fuse.DT_File}

	end := int(len(s.itemz))
	if end >= 100 {
		end = 100
	}

	for i := 0; i < end; i++ {
		dirz[i+foffset].Inode = uint64(i + inodeoffset)
		dirz[i+foffset].Name = fmt.Sprintf("%02d", i)
		dirz[i+foffset].Type = fuse.DT_Dir
	}
	sdirs := dirz[0 : end+foffset]

	return sdirs, nil
}

func (tappertrackernode) Open(req *fuse.OpenRequest, resp *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	//TODO: open self here
	return fs.DataHandle([]byte("")), nil
}

//ok
func (tappertrackernode) Attr() fuse.Attr {
	a := generic_attr()
	//TODO: report self size here
	a.Inode = 2
	a.Size = 0
	a.Blocks = (0 / 512)
	a.Mode = 0555
	a.Nlink = 1 // correct?//FIXME
	return a
}

//ok
func (s tapperdirnode) Attr() fuse.Attr {
	a := generic_attr()

	a.Inode = s.i + inodeoffset
	a.Size = 4096
	a.Blocks = 8
	a.Mode = os.ModeDir | 0555
	a.Nlink = 2 // correct?//FIXME
	return a
}

func (s tapperdirnode) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var foobar []fuse.Dirent

	ibase := (s.i << 18) + 128

	for i := range s.itemz {
		item := fuse.Dirent{Inode: uint64(i) + ibase, Name: s.itemz[i], Type: fuse.DT_Link}
		foobar = append(foobar, item)
	}
	return foobar, nil
}

func (s tapperdirnode) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {

	ibase := (s.i << 18) + 128

	// TODO: binary search
	for i := range s.itemz {
		if name == s.itemz[i] {
			return tapperbinlink{inode: ibase + uint64(i)}, nil
		}
	}

	return nil, fuse.ENOENT
}

//important: do not add Getattr, it will not work

//ok
func (s tapperbinlink) Attr() fuse.Attr {
	a := generic_attr()

	a.Inode = s.inode
	a.Size = 10
	a.Blocks = 8
	a.Mode = 0555 | os.ModeSymlink
	a.Nlink = 1
	return a
}

func (tapperbinlink) Readlink(req *fuse.ReadlinkRequest, intr fs.Intr) (string, fuse.Error) {
	return "../tracker", nil
}
