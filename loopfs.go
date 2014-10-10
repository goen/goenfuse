package main

import (
	"fmt"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
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
