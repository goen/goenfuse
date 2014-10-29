// Hellofs implements a simple "hello world" file system.
package main

import (
	"fmt"

	"encoding/gob"
	"bufio"

	"os"
	"os/signal"
	"sync"
)


func (d dump) write(op Fileop) {
	d.Lock()
	defer d.Unlock()

	if d.t != nil {
		d.enc.Encode(op)
	}
}

type dump struct {
	sync.Mutex
	t *bufio.Writer
	enc *gob.Encoder
}

func main() {
	var d dump

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

	// open the writer
	s, errr := os.OpenFile(mpoint_gbin+"/loop", os.O_WRONLY, 0200)
	if errr != nil {
		fmt.Println(errr)
		return
	}

	d.t = bufio.NewWriter(s)
	d.enc = gob.NewEncoder(d.t)

	for !bin.u || !loop.u {

		//wait for signal
		for sig := range sigchan {
			fmt.Println("stopped!", sig)
			break
		}

		d.write(Fileop{Code: OP_UMOUNT, File: ""})
		d.Lock()
		if d.t != nil {
			d.t.Flush()
			d.t = nil
			s.Close()
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
