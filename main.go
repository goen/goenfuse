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
	f.u = false
	return f, e
}

func (f *ffs) umount() (err error) {
	if f.u {
		return nil
	}
	// taken from the fs/fstestutil/mounted.go
	for tries := 0; tries < 100; tries++ {
		err = fuse.Unmount(f.dir)
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		f.u = true
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
	u    bool //umounted ok
}

func scan_path(p string) (items []string) {
	fmt.Println("Scanning path ", p)

	filepath.Walk(p, func(path string, f os.FileInfo, _ error) error {

		if p == path {
			return nil
		}

		//	fmt.Println("|>>| ",  path)

		base := filepath.Base(path)

		items = append(items, base)
		return nil
	})

	//	for i := range items {
	//	fmt.Println("|| ", items[i])
	//	}

	return items
}

func main() {
	path := []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin", "/usr/bin", "/sbin", "/bin"}
	pitems := [][]string{}

	for i := range path {
		pitems = append(pitems, scan_path(path[i]))
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

	go loop.try_serve(looperfs{})
	go bin.try_serve(tapperfs{r: tapperrootnode{dirs: uint64(len(path)), itemz: pitems}})

	//wait until mounted
	loop.check_err()
	bin.check_err()

	for !bin.u || !loop.u {

		//wait for signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		for sig := range c {
			fmt.Println("quitting!", sig)
			break
		}

		if loop.umount() != nil {
			fmt.Println("Umounting ", loop.dir, " failed")
		}
		if bin.umount() != nil {
			fmt.Println("Umounting ", bin.dir, " failed")
		}

		if !bin.u || !loop.u {
			fmt.Println("Please, stop using & quit the drive")
			fmt.Println("Then, try again..")
		}

	}
}
