//this is the common fuse glue
package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

/*
the order is: mount, putcontext, serve, check, umount, destroy
*/

// workaround for mkdir panic
func failsafe_mkdir_all(dir string, perm os.FileMode) error {
	var wg sync.WaitGroup
	wg.Add(1)

	panik := true

	go func(dir string, panik *bool) {
		defer wg.Done()
		if os.MkdirAll(dir, perm) == nil {
			*panik = false
		}
	}(dir, &panik)

	wg.Wait()
	if panik {
		return fmt.Errorf("Failsafe make directory failed.")
	}
	return nil
}

func mount(dir string) (f Ffs, e error) {
	_, e = os.Stat(f.dir)
	f.lack = e != nil
	f.dir = dir
	f.u = false
	if f.lack {
		//	e = os.MkdirAll(dir, 755)	//this may panic
		e = failsafe_mkdir_all(dir, 755)
		if e != nil {
			return f, e
		}
		fmt.Println("maked ", dir)
	}

	f.monut()

	return f, e
}

func (f *Ffs) umount() (err error) {
	if f.u {
		return nil
	}
	// taken from the fs/fstestutil/mounted.go
	for tries := 0; tries < 100; tries++ {

		err := f.unmount()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		f.u = true
		return nil
	}
	return err
}

func destroy(f Ffs) {
	destory(f)
	if f.lack {
		os.RemoveAll(f.dir)
	}
}

// my fuse fs
type Ffs struct {
	dir   string
	lack  bool
	be    fbackend
	u     bool //umounted ok
	stuff stuffer
}
