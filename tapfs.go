package main

import (
	"fmt"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// TapFS implements the tap part of the fs. Used for listening
// for file access
type TapFS struct{}

func (TapFS) Root() (fs.Node, fuse.Error) {
	if debug {
		fmt.Println("TAPFS::ROOT")
	}
	return Dir{}, nil
}
