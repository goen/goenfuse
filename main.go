// Hellofs implements a simple "hello world" file system.
package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

const (
	coolflag = "override-my-path"
)

var coolflagg = flag.String(coolflag, "", "Skip binary self-path lookup and self-check, by  ")

// workaround

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

//begin ffs

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
	stuff interface{}
}

//end ffs

func scan_path(p string) (items []string, has_me bool) {

	filepath.Walk(p, func(path string, f os.FileInfo, _ error) error {

		if p == path {
			return nil
		}

		base := filepath.Base(path)

		if base == self_file {
			has_me = true
		}

		items = append(items, base)
		return nil
	})

	return items, has_me
}

func tracker_main() {
	fmt.Println("HELLO FROM TRACKER")

	// TODO: dump ENV
	fmt.Println("ENV:", os.Environ())

	// TODO: dump ARGS
	fmt.Println("ARGS:", os.Args)

	// TODO: dump ARGS
	fmt.Println("EXEC:run the actual binary:", os.Args)

	/*
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Waiting for command to finish...")
		err = cmd.Wait()
		fmt.Println("Command finished with error: ", err)
	*/
	// TODO: wait
	fmt.Println("TODO:wait for the actual binary to complete")
	/**/
	fmt.Println("BYE BYE FROM TRACKER")
}

func main() {
	flag.Parse()
	var myself self
	myself.set(*coolflagg)

	var tracker bool

	if filepath.Base(os.Args[0]) != self_file {
		fmt.Println("003")
		dir := self_2digit_dir()
		if dir != 255 {
			tracker = true
		} else {
			tracker = false
		}
	}

	if len(os.Args) > 2 {
		fmt.Println("002")
		tracker = true
	}

	if myself.is() {
		fmt.Println("001")
		tracker = false
	}

	if tracker {
		fmt.Println("004")
		tracker_main()
		return
	}

	// this is the path

	// FIXME: load here actual env from path
	var path []string = self_path

	// load the path items from path

	var pitems_array [100][]string
	pitems := pitems_array[0:len(path)]
	_ = pitems

	var mybinwhere uint32 = 0xffff

	var wg sync.WaitGroup
	for i := range path {
		wg.Add(1)
		go func(j uint32) {
			var has_me bool
			pitems_array[j], has_me = scan_path(path[j])
			if has_me {
				mybinwhere = j
			}
			wg.Done()
		}(uint32(i))
	}
	wg.Wait()
	//done loading path items

	self_locs := []string{}

	if !myself.is() {

		// look at my binary in path

		if mybinwhere < uint32(len(path)) {
			p := path[mybinwhere]
			self_locs = append(self_locs, p+"/"+self_file, p+"/"+os.Args[0])
		}

		// this binary may be in this dir
		pwd, err := os.Getwd()
		if err == nil {
			self_locs = append(self_locs, pwd+"/"+self_file, pwd+"/"+os.Args[0])
		}

		// consider add lookup by readlink -f /proc/$pid/exe

		// check binary contains the magic string selfrtg
		for i := range self_locs {
			if self_check(self_locs[i]) {
				myself.set(self_locs[i])
				break
			}
		}

		if !myself.is() {
			fmt.Println("The `" + self_file + "` file not found.\n" +
				"Run " + self_file + " --" + coolflag + "=/../.." + self_file)
			return
		}
	}

	//capturing signals before and after mount
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	loop, errl := mount(mpoint_gloop)
	bin, errb := mount(mpoint_gbin)

	if errl == nil {
		errl = loop.mount()
	}
	if errb == nil {
		errb = bin.mount()
	}

	if errl != nil || errb != nil {
		fmt.Println("Mount failed: ", errl)
		fmt.Println("Already mounted or stale mount")
		if errl == nil {
			destroy(loop)
		}
		if errb == nil {
			destroy(bin)
		}
		return
	}
	defer destroy(loop)
	defer destroy(bin)

	//bazil specific
	loop.stuff = looperfs{}
	bin.stuff = tapperfs{r: tapperrootnode{itemz: pitems, s: &myself}}

	go loop.serve()
	go bin.serve()

	//wait until mounted
	loop.check()
	bin.check()

	for !bin.u || !loop.u {

		//wait for signal
		for sig := range sigchan {
			fmt.Println("stopped!", sig)
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
