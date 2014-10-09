package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// TapFS implements the tap part of the fs. Used for listening
type TapFS struct{}

func (TapFS) Root() (fs.Node, fuse.Error) {
	return Dir{}, nil
}
