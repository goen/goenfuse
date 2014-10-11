package main

import (
	"fmt"
	"os"
)

func tap(bin string, tap string) {
	// clean stale tap
	_ = os.Remove(bin + "\n")
	{
		err := os.Rename(bin, bin+"\n")
		if err != nil {
			fmt.Println("Can't putaway")
		}
	}
	// TODO: consider soft-linking
	{
		err := os.Link(tap, bin)
		if err != nil {
			fmt.Println("Can't hardlink")
		}
	}
}

func untap(bin string) {
	{
		err := os.Remove(bin)
		if err != nil {
			fmt.Println("Can't remove")
		}
	}
	{
		err := os.Rename(bin+"\n", bin)
		if err != nil {
			fmt.Println("Can't putback")
		}
	}
}

func main() {

	untap("bin/tapper")
	//	fmt.Println("Hello world")

	//	tap("bin/tapper", "../tracker/tracker")
}
