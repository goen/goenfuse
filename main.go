// Hellofs implements a simple "hello world" file system.
package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"bitbucket.org/kardianos/osext"
	"os/exec"
	"strings"
	"syscall"
//	"time"
//	"io/ioutil"
//	"io"
)

const (
	coolflag = "override-my-path"
)

var coolflagg = flag.String(coolflag, "", "Skip binary self-path lookup and self-check, by  ")

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

func tracker_main() int {
	trynotify("EXEC //" + filepath.Clean(os.Args[0])+"//")

	// clean the PATH

	var newpath []string
	p := strings.Split(os.Getenv("PATH"), ":")
	for i := range p {
		pi := filepath.Clean(p[i])
		if len(pi) >= 3  {
			if decodedir(pi[len(pi)-2:]) != 255 {
				continue
			}
		}

		newpath = append(newpath, p[i])
	}
	os.Setenv("PATH", strings.Join(newpath, ":"))

	xec := filepath.Clean(selffile("abspaths")[underscore_hack()] + "/" + filepath.Base(os.Args[0]))

	cmd := exec.Command(xec, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("err1")
		return -1
	}

	err = cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		fmt.Println("err2")
		return -2
	}

	return 0
}

func trynotify(s string) bool {
	// notify the listener
	pipe, err := tapopen()
	if err == nil {
		fmt.Fprintln(pipe, s)
		pipe.Close()
		return true
	}
	return false
}

func main() {
	var tracker bool

	if filepath.Base(os.Args[0]) != self_file {
		tracker = (self_2digit_dir() != 255)
	}

	if len(os.Args) > 2 {
		tracker = true
	}

	if selfer() {
		tracker = true
	}

	if selfish_arg() {
		tracker = false
	}

	if tracker {
		os.Exit(tracker_main())
	}

	// welcome to the fuse part
	flag.Parse()
	var myself self
	myself.set(*coolflagg)

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

		// lookup by readlink -f /proc/$pid/exe
		myloc, err2 := osext.Executable()
		if err2 != nil {
			self_locs = append(self_locs, myloc)
		}

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

	loop.stuff = loopcontext()
	bin.stuff = tapcontext(pitems, &myself, path)

	if errl == nil {
		errl = loop.putcontext()
	}
	if errb == nil {
		errb = bin.putcontext()
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

		trynotify("UMOUNTED")

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
