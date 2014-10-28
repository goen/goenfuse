// Hellofs implements a simple "hello world" file system.
package main

import (
	"fmt"

	"os"
	"os/signal"
)

func main() {
	//capturing signals before and after mount
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	loop, errl := mount(mpoint_gloop)
	bin, errb := mount(mpoint_gbin)

	loop.stuff = loopcontext()
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
