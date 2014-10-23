//this is the conversion from linux stat to the fuse attr
package main

import (
	"os"
	"syscall"
	"time"
)

type Stat syscall.Stat_t

func StatSys(fi os.FileInfo) Stat {
	a := fi.Sys()
	if a == nil {
		return Stat{}
	}
	b := a.(*syscall.Stat_t) //whatever

	var c Stat = Stat(*b)

	return c
}

func (s Stat) Inode() uint64 { // inode number
	return s.Ino
}
func (s Stat) Size_() uint64 { // size in bytes
	return uint64(s.Size)
}
func (s Stat) Blocks_() uint64 { // size in blocks
	return uint64(s.Blocks)
}
func (s Stat) Atime() time.Time { // time of last access
	return time.Unix(s.Atim.Sec, s.Atim.Nsec)
}
func (s Stat) Mtime() time.Time { // time of last modification
	return time.Unix(s.Mtim.Sec, s.Mtim.Nsec)
}
func (s Stat) Ctime() time.Time { // time of last inode change
	return time.Unix(s.Ctim.Sec, s.Ctim.Nsec)
}
func (s Stat) Crtime() time.Time { // time of creation (OS X only)
	return time.Unix(0, 0)
}

func (s Stat) Mode_() os.FileMode { // file mode
	return os.FileMode(s.Mode)
}
func (s Stat) Nlink_() uint32 { // number of links
	// FIXME: why bazil fuse nlink is only 32bit? bug?
	return uint32(s.Nlink)
}

func (s Stat) Rdev_() uint32 { // device numbers
	// FIXME: why bazil fuse nlink is only 32bit? bug?
	return uint32(s.Rdev)
}

func (s Stat) Flags() uint32 { // chflags(2) flags (OS X only){
	return 0
}
