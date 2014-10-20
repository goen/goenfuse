package main

import (
	//	"fmt"
	"bufio"
	"bytes"
	"os"
	"sync"
)

const (
	self_file = "goenfuse"
	selfrtg   = "\xb8\x4c\x10\x44\x00\x00\x00\x00\x00\x00\x8b\x55\x84\xdb\xde"
)

type self struct {
	sync.Mutex
	path string
	size uint64
}

func (s self) get() (string, uint64) {
	s.Lock()
	defer s.Unlock()
	return s.path, s.size
}

func (s self) is() bool {
	s.Lock()
	defer s.Unlock()
	return s.path != ""
}

func (s *self) set(str string) {
	s.Lock()
	defer s.Unlock()
	s.path = str
	fi, err := os.Lstat(s.path)
	if err == nil {
		s.size = uint64(fi.Size())
	}
}

func self_check(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	tag := []byte(selfrtg)
	br := bufio.NewReader(file)

	err = bufio.ErrBufferFull
	line := []byte{}

	for {
		line, err = br.ReadSlice(tag[0])
		if err == bufio.ErrBufferFull {
			continue
		}
		if err != nil {
			return false
		}
		line, err = br.ReadSlice(tag[len(tag)-1])
		if err == bufio.ErrBufferFull {
			//			fmt.Println("LEN IS ENOUGH :D", len(line))
			continue
		} else if err != nil {
			return false
		}
		if len(line)+1 != len(tag) {
			continue
		}

		if bytes.Equal(line, tag[1:]) {
			//			fmt.Println("FOUND!!!")
			return true
		}
	}
	return false
}
