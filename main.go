// Hellofs implements a simple "hello world" file system.
package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

//tapdir
type tapdir struct{}

func (tapdir) Attr() fuse.Attr {
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
}

func (tapdir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if name == "hi" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

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
	{Inode: 3, Name: "world", Type: fuse.DT_File},
}

func (Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	if debug {
		fmt.Println("DIR::ReadDir")
	}

	return dirDirs, nil
}

// tapfile implements an executable asset used for tracking the activity
type tapfile struct{}

func (tapfile) Attr() fuse.Attr {
	return fuse.Attr{Inode: 2, Mode: 0555, Size: uint64(1834576)}
}

func (tapfile) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	data, _ := Asset("tracker")
	return data, nil
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

//begin ffs stuff

func mount(dir string) (f ffs, e error) {
	f.dir = dir

	_, e = os.Stat(f.dir)
	f.lack = e != nil

	if f.lack {
		err := os.MkdirAll(f.dir, 755)
		if err != nil {
			return f, err
		}
	}

	f.c, e = fuse.Mount(f.dir)
	return f, e
}

func (f ffs) umount() (err error) {
	// taken from the fs/fstestutil/mounted.go
	for tries := 0; tries < 1000; tries++ {
		err = fuse.Unmount(f.dir)
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		return nil
	}
	return err
}

func (f ffs) try_serve(s fs.FS) {
	fs.Serve(f.c, s)
}

func (f ffs) check_err() error {
	<-f.c.Ready
	if err := f.c.MountError; err != nil {
		return err
	}
	return nil
}

func destroy(f ffs) {
	f.c.Close()

	if f.lack {
		os.RemoveAll(f.dir)
	}
}

// my fuse fs
type ffs struct {
	dir  string
	lack bool
	c    *fuse.Conn
}

//end ffs stuff

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func scan_path(p string) error {
	fmt.Println("Scanning path ", p)

	return filepath.Walk(p, visit)
}

func main() {
	path := []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin", "/usr/bin", "/sbin", "/bin"}

	for i := range path {
		scan_path(path[i])
	}

	loop, errl := mount("goenloop")
	bin, errb := mount("goenbin")

	if errl != nil {
		panic(errl)
	}
	defer destroy(loop)
	if errb != nil {
		panic(errb)
	}
	defer destroy(bin)

	go loop.try_serve(LoopFS{})
	go bin.try_serve(TapFS{})

	//wait until mounted
	loop.check_err()
	bin.check_err()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		fmt.Println("quitting!", sig)
		break
	}

	loop.umount()
	bin.umount()
}
