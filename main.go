// Hellofs implements a simple "hello world" file system.
package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"log"
	"os"
	"time"
	"fmt"
)

// Dir implements both Node and Handle for the root directory.
type Dir struct{}

func (Dir) Attr() fuse.Attr {
	if debug {
		fmt.Println("DIR::ATTR")
	}

	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
}
func (Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if debug {
		fmt.Println("DIR::LOOKUP:",name)
	}

	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
}

func (Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	if debug {
		fmt.Println("DIR::ReadDir")
	}

	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

const greeting = "hello, world\n"

func (File) Attr() fuse.Attr {
	if debug {
		fmt.Println("FILE::attr")
	}

	return fuse.Attr{Inode: 2, Mode: 0444, Size: uint64(len(greeting))}
}
func (File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	if debug {
		fmt.Println("FILE::readall")
	}

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
	trackmount := "remapped"

	os.MkdirAll(goenmount, 755)

	// mount the mnt
	c1, err := fuse.Mount(goenmount)
	if err != nil {
	if debug {
		log.Fatal(err)
	}}

	c2, err := fuse.Mount(trackmount)
	if err != nil {
	if debug {
		log.Fatal("alpha",err)
	}}

	//shedule deletion
	go mountdeleter(goenmount)
	go mountdeleter(trackmount)

	defer c1.Close()
	defer c2.Close()

	err = fs.Serve(c1, LoopFS{})
	if err != nil {
	if debug {
		log.Fatal("bb",err)
	}}

	err = fs.Serve(c2, TapFS{})
	if err != nil {
	if debug {
		log.Fatal("ccc",err)
	}}

	<-c1.Ready
	<-c2.Ready

	if err := c1.MountError; err != nil {
	if debug {
		log.Fatal("dddd",err)
	}}

	if err := c2.MountError; err != nil {
	if debug {
		log.Fatal("eeee",err)
	}}

}
