// Hellofs implements a simple "hello world" file system.
package main

import (
	"fmt"

	"encoding/gob"
	"bufio"

	"os"
	"os/signal"
	"sync"
	"time"
)


func (d *dump) write(op Fileop) {
	d.Lock()
	defer d.Unlock()

	if d.t != nil {
		d.enc.Encode(op)
	}
}

type dump struct {
	sync.Mutex
	tt byte
	terminate bool
	s *os.File
	t *bufio.Writer
	enc *gob.Encoder
}

func busy(d *dump) {
	d.Lock()
	defer d.Unlock()

	for d.tt == 0 {
		d.Unlock()
		time.Sleep(1000000)
		d.Lock()
	}
}

func pipetoucher(d *dump) {
	d.Lock()
	defer d.Unlock()

	if d.tt != 0 || d.terminate {
		return
	}

	d.terminate = true

	s, err := os.Open(mpoint_gbin+"/loop")
	if err != nil {
		fmt.Println(err)
		return
	}

	s.Close()
}

func pipeopener(d *dump) {

	// open the writer
	s, errr := os.OpenFile(mpoint_gbin+"/loop", os.O_WRONLY, 0200)
	if errr != nil {
		fmt.Println(errr)
		return
	}

	d.Lock()
	defer d.Unlock()

	d.tt = 1
	d.s = s

	if d.terminate {
		return
	}

	d.t = bufio.NewWriter(d.s)
	d.enc = gob.NewEncoder(d.t)
}

func main() {
	var d dump
	d.t = nil

	//capturing signals before and after mount
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	loop, errl := mount(mpoint_gloop)
	bin, errb := mount(mpoint_gbin)

	loop.stuff = loopcontext(&d)
	bin.stuff = tapcontext()

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

	fmt.Println("Waiting for the dump")
	go pipeopener(&d)

	for !bin.u || !loop.u {

		//wait for signal
		for sig := range sigchan {
			fmt.Println("stopped!", sig)
			break
		}

		d.write(Fileop{Code: OP_UMOUNT, File: ""})

		pipetoucher(&d)
		busy(&d)
		d.Lock()
		if d.tt == 1 {
			if d.t != nil {
				d.t.Flush()
			}
			d.tt = 2
			d.t = nil
			d.s.Close()
		}
		d.Unlock()

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
