package main

import (

)

const (
	OP_HELLO = iota
	OP_UMOUNT
	OP_GETATTR
	OP_OPENDIR
	OP_OPEN
	OP_CHMOD
	OP_CHOWN
	OP_TRUNC
	OP_UTIMENS
	OP_READLNK
	OP_MKNOD
	OP_MKDIR
	OP_UNLINK
	OP_RMDIR
	OP_SLINK
	OP_RENAME
	OP_LINK
	OP_ACCESS
	OP_CREATE
	OP_READ
	OP_WRITE
	OP_RELEASE
	OP_FLUSH
	OP_FSYNC
)

type Fileop struct {
	Code uint8
	Openid uint64
	Errno int32
	Nsec       int64
	File string
	Extra string
}


