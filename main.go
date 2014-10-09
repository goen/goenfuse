// Hellofs implements a simple "hello world" file system.
package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"log"
	"os"
	"time"
)

// Dir implements both Node and Handle for the root directory.
type Dir struct{}

func (Dir) Attr() fuse.Attr {
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
}
func (Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
}

func (Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

const greeting = "hello, world\n"

func (File) Attr() fuse.Attr {
	return fuse.Attr{Inode: 2, Mode: 0444, Size: uint64(len(greeting))}
}
func (File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	return []byte(greeting), nil
}

func mountdeleter(p string) {
	time.Sleep(10 * time.Second)

	mountkiller(p)
	os.RemoveAll(p)
}

func mountkiller(p string) {
	err := killer(p)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func killer(p string) (err error) {
	for tries := 0; tries < 1000; tries++ {
		err = fuse.Unmount(p)
		if err != nil {

			time.Sleep(10 * time.Millisecond)
			continue
		}
		return nil
	}
	return err
}

func main() {

	goenmount := "goenmount"

	os.MkdirAll(goenmount, 755)

	c, err := fuse.Mount(goenmount)
	if err != nil {
		log.Fatal(err)
	}

	go mountdeleter(goenmount)

	defer c.Close()

	err = fs.Serve(c, LoopFS{})
	if err != nil {
		log.Fatal(err)
	}

	<-c.Ready

	if err := c.MountError; err != nil {
		log.Fatal(err)
	}

}
