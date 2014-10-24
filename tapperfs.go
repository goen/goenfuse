// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type tapper_root struct {
	nodefs.Node
}
