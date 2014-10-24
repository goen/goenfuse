// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

func loopcontext() nodefs.Node {
	return looper_root{}
}

type looper_root struct {
	nodefs.Node
}
