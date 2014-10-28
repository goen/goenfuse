//this is the self checking code
package main

import (
	"bitbucket.org/kardianos/osext"

	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"io/ioutil"
	"strings"
	"fmt"
	"io"
)

const (
	self_file = "goenfuse"
	selfrtg   = "\xb8\x4c\x10\x44\x00\x00\x00\x00\x00\x00\x8b\x55\x84\xdb\xde"
)

var self_path = []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin",
	"/usr/bin", "/sbin", "/bin"}

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
	if !selfcheck {
		return true
	}
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

			continue
		} else if err != nil {
			return false
		}

		if len(line) >= len(tag) {
			if bytes.Equal(line[len(line)-len(tag):], tag) {
				return true
			}
		}

		if len(line)+1 != len(tag) {
			continue
		}

		if bytes.Equal(line, tag[1:]) {
			return true
		}
	}
	return false
}

func selfish_arg() bool {
	if len(os.Args) != 2  {
		return false
	}

	g := os.Args[1]
	l := len(coolflag) + 2

	if len(g) < l {
		return false
	}

	return bytes.Equal([]byte(g[0:l]), []byte("--" + coolflag))
}


type tee struct {
	sync.Mutex
	stop bool
	io.ReadCloser
	io.WriteCloser
}

func (t tee) Close() error {
	t.ReadCloser.Close()
	t.WriteCloser.Close()
	return nil
}

func (t tee) Kill() {
	t.stop = true
}

func teeopen() (rdc tee, err error) {
	a, b := os.Open(mpoint_gbin+"/write")
	if b != nil {
		return rdc, fmt.Errorf("Unable to open the input pipe:",b)
	}

	out, errr := tapopen(false)
	if errr != nil {
		a.Close()
		return rdc, fmt.Errorf("Error: opening out pipe:", errr)
	}

	return tee{ReadCloser:a, WriteCloser:out}, nil
}

func tapopen( who bool) (rdr io.WriteCloser, err error){
	myloc, err2 := osext.Executable()
	if err2 != nil {
		return nil, fmt.Errorf("Unable to find the output pipe:",err2)
	}
	var pipe string
	if who {
		pipe = "/write"
	} else {
		pipe = "/read"
	}

	a, b := os.Create(filepath.Dir(myloc)+pipe)
	if b != nil {
		return nil, fmt.Errorf("Unable to open the output pipe")
	}

	return io.WriteCloser(a), nil
}

func selffile( file string) []string {
	myloc, err2 := osext.Executable()
	if err2 != nil {
		return []string{}
	}
	content, err := ioutil.ReadFile(filepath.Dir(myloc)+"/"+file)
	if err != nil {
		return []string{}
	}
	return strings.Split(string(content), "\n")
}

func selfer() bool {
	myloc, err2 := osext.Executable()
	if err2 != nil {
		return false
	}

	dir := filepath.Base(filepath.Dir(myloc))

	return dir == mpoint_gbin
}

func decodedir(d string) uint8 {
	dir := []byte(d)
	if len(dir) != 2 {
		return 255
	}
	if dir[0] < '0' || dir[0] > '9' || dir[1] < '0' || dir[1] > '9' {
		return 255
	}

	num := 10*(dir[0]-'0') + (dir[1] - '0')
	if int(num) > len(self_path) {
		return 255
	}
	return num
}

func underscore_hack() uint8 {
	p := os.Getenv("_")
	dir := filepath.Base(filepath.Dir(p))
	return decodedir(dir)
}

func self_2digit_dir() uint8 {
	// lookup by readlink -f /proc/$pid/exe
	myloc, err2 := osext.Executable()
	if err2 != nil {
		return 255
	}

	dir := filepath.Base(filepath.Dir(myloc))
	return decodedir(dir)
}
