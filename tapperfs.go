// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"fmt"
	"os"
)

func tapcontext(i [][]string, z *self) nodefs.Node {
	return tapper_root{itemz: i, self: z}
}

type tapper_root struct {
	nodefs.Node
	itemz [][]string // len(itemz) = 1 + maximum name
	*self
}

func (r tapper_root) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	var dirz [103]fuse.DirEntry
	foffset := 3

	dirz[0] = fuse.DirEntry{Name: "tracker", Mode: 0555}
	dirz[1] = fuse.DirEntry{Name: ".track", Mode: 0555}
	dirz[2] = fuse.DirEntry{Name: ".untrack", Mode: 0555}

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
