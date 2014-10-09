package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// LoopFS implements the loop part of the fs.
type LoopFS struct{}

func (LoopFS) Root() (fs.Node, fuse.Error) {
	return Dir{}, nil
}
