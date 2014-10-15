package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"fmt"
	"os"
)

// LoopFS implements the loop part of the fs.
// A FUSE filesystem that shunts all request to an underlying file
// system.
type LoopFS struct{}

func (LoopFS) Root() (fs.Node, fuse.Error) {
	if debug {
		fmt.Println("LOOPFS::ROOT")
	}
	return Dir{}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct{}

func (Dir) Attr() fuse.Attr {
	if debug {
		fmt.Println("DIR::ATTR")
	}

	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
}
func (Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if debug {
		fmt.Println("DIR::LOOKUP:", name)
	}

	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
	{Inode: 3, Name: "world", Type: fuse.DT_File},
}

func (Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	if debug {
		//		fmt.Println("DIR::ReadDir")
	}

	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

const greeting = "hello, world\n"

func (File) Attr() fuse.Attr {
	if debug {
		fmt.Println("FILE::attr")
	}

	return fuse.Attr{Inode: 2, Mode: 0444, Size: uint64(len(greeting))}
}
func (File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	if debug {
		fmt.Println("FILE::readall")
	}

	return []byte(greeting), nil
}
