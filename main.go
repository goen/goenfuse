// Hellofs implements a simple "hello world" file system.
package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"fmt"
	"log"
	"os"
	"time"
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
		fmt.Println("DIR::LOOKUP:", name)
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

func umount(p string) (err error) {
	// taken from the fs/fstestutil/mounted.go
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
	os.MkdirAll(trackmount, 755)

	// mount the mnt
	c1, err1 := fuse.Mount(goenmount)
	if err1 != nil {
		if debug {
			log.Fatal(err1)
		}
	}
	defer c1.Close()

	c2, err2 := fuse.Mount(trackmount)
	if err2 != nil {
		if debug {
			log.Fatal(err2)
		}
	}
	defer c2.Close()

	go fs.Serve(c1, LoopFS{})
	go fs.Serve(c2, TapFS{})

	<-c1.Ready
	<-c2.Ready

	if err := c1.MountError; err != nil {
		if debug {
			log.Fatal("dddd", err)
		}
	}

	if err := c2.MountError; err != nil {
		if debug {
			log.Fatal("eeee", err)
		}
	}

	//shedule deletion
	fmt.Println("sleeping")
	time.Sleep(6 * time.Second)
	fmt.Println("umounting")
	_ = umount(goenmount)
	_ = umount(trackmount)

	os.RemoveAll(goenmount)
	os.RemoveAll(trackmount)

}
