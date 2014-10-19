package main

import (
	"fmt"
	"os"
	//	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("HELLO FROM TRACKER")
	// TODO: dump ENV
	fmt.Println("TODO:dump ENV")
	// TODO: dump ARGS
	fmt.Println("TODO:dump ARGS")
	// TODO: dump ARGS
	fmt.Println("TODO:run the actual binary:", os.Args)

	if filepath.Base(os.Args[0]) == "tracker" {
		return
	}
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
