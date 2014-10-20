package main

import (
	//	"fmt"
	"bufio"
	"bytes"
	"os"
)

const (
	self_file = "goenfuse"
	selfrtg   = "\xb8\x4c\x10\x44\x00\x00\x00\x00\x00\x00\x8b\x55\x84\xdb\xde"
)

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
