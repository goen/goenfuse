package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"fmt"
	"os"
	"time"
)

//ok
func generic_attr(now bool, l bool) fuse.Attr {
	u := time.Unix(0, 0)
	if now {
		u = time.Now()
	}
	var d uint32 = 0xB1F5     //BIn-FS
	var f uint32 = 0x00121012 //don't modify
	if l {
		d = 0x70F5 //LOop-FS
		f = 0      //modify
	}

	return fuse.Attr{
		Atime: u, Mtime: u, Ctime: u, Crtime: u,
		Uid:   uint32(os.Geteuid()),
		Gid:   uint32(os.Getegid()),
		Rdev:  d,
		Flags: f,
	}
}

// tapperdir implements both Node and Handle for the root directory.
//HandleReader, HandleReadDirer

// tapperFS dirs: root dir & the various dirs
//

//ok
type tapperfs struct {
	r tapperrootnode
}

//ok
type tapperrootnode struct {
	dirs  uint64 // = 1 + maximum name
	itemz [][]string
}

//ok
type tappertrackernode struct {
	//???
}

//ok
type tapperdirnode struct {
	i uint64 //name = i, inode = i + 3
	// pointer to slice of strings

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
	a := generic_attr(false, false)
	a.Inode = 1
	a.Size = 4096
	a.Blocks = 8
	a.Mode = os.ModeDir | 0555
	a.Nlink = 8 // correct?//FIXME
	return a
}

//ok
func (s tapperrootnode) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if name == "tracker" {
		return tappertrackernode{}, nil
	}

	var i int

	n, err := fmt.Sscanf(name, "%02d", &i)
	if (err != nil) || (n != 1) || (uint64(i) >= s.dirs) {
		return nil, fuse.ENOENT
	}

	return tapperdirnode{i: uint64(i), itemz: s.itemz[i]}, nil
}

//ok
func (s tapperrootnode) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var dirz [101]fuse.Dirent

	dirz[0] = fuse.Dirent{Inode: 2, Name: "tracker", Type: fuse.DT_File}

	end := int(s.dirs)
	if end >= 100 {
		end = 100
	}

	for i := 0; i < end; i++ {
		dirz[i+1].Inode = uint64(i + 3)
		dirz[i+1].Name = fmt.Sprintf("%02d", i)
		dirz[i+1].Type = fuse.DT_Dir
	}
	sdirs := dirz[0 : end+1]

	return sdirs, nil
}

func (tappertrackernode) Open(req *fuse.OpenRequest, resp *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	//fmt.Print("Open Called: ", req, resp, intr)

	a, err := Asset("tracker")
	if err != nil {
		return nil, nil //FIXME return error here
	}

	return fs.DataHandle(a), nil
}

//ok
func (tappertrackernode) Attr() fuse.Attr {
	a := generic_attr(false, false)

	a.Inode = 2
	a.Size = bin_tracker_size
	a.Blocks = (bin_tracker_size / 512)
	a.Mode = 0555
	a.Nlink = 1 // correct?//FIXME
	return a
}

//ok
func (s tapperdirnode) Attr() fuse.Attr {
	a := generic_attr(false, false)

	a.Inode = s.i + 3
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

		//		fmt.Println("reading dir:", s.itemz[i])
		item := fuse.Dirent{Inode: uint64(i) + ibase, Name: s.itemz[i], Type: fuse.DT_Link}
		foobar = append(foobar, item)
	}
	/*
		item1 := fuse.Dirent{Inode: 1337, Name: "hello", Type: fuse.DT_Link}
		item2 := fuse.Dirent{Inode: 1338, Name: "world", Type: fuse.DT_Link}
		foobar = append(foobar, item1, item2)
	*/
	return foobar, nil
}

func (s tapperdirnode) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {

	ibase := (s.i << 18) + 128

	// TODO: binary search
	for i := range s.itemz {
		//		fmt.Println("||", s.itemz[i])
		if name == s.itemz[i] {
			return tapperbinlink{inode: ibase + uint64(i)}, nil
		}
	}
	/*
		if name == "hello" {
			return tapperbinlink{inode: 1337}, nil
		}
		if name == "world" {
			return tapperbinlink{inode: 1338}, nil
		}

		fmt.Println("Lookup!!!!!!!1", name, intr)
	*/
	return nil, fuse.ENOENT
}

//important: do not add Getattr, it will not work

//ok
func (s tapperbinlink) Attr() fuse.Attr {
	a := generic_attr(false, false)

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
