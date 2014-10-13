package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
	TODO: compile the binary into a binary.
	TODO: when I run tap i check if they're already tapped
	TODO: when untap check if they arent really tapped
*/

func do_tapbin(dir string) int {

	data, _ := Asset("tracker")
	if len(data) == 0 {
		return -1
	}
	ioutil.WriteFile(dir+"/\n", data, 0755) // -rwxr-xr-x
	return 0
}

func rm_tapbin(dir string) error {
	return os.Remove(dir + "/\n")
}

func is_tapdir(dir string) {

}

// the purpose of this is to check a binary and see if it's a tapbin. e.g. my wrapper binary
func is_tapbin(bin string) int {
	if strings.HasSuffix(bin, "\n") {
		fmt.Println("this is a backup of a binary. unless you ran tap twice wich shouldnt happen.")
	}
	fmt.Println("this probably a tapping binary.")
	fmt.Println("but could be a orig binary.")
	return 1337
}

func dotap(bin string) int {
	tap := filepath.Dir(bin) + "/\n"

	// clean stale tap
	_ = os.Remove(bin + "\n")
	{
		err := os.Rename(bin, bin+"\n")
		if err != nil {
			fmt.Println("Can't put away")
			return -1
		}
	}
	if tapfile_hardlink {
		err := os.Link(tap, bin)
		if err == nil {
			return 0
		}
	}
	if tapfile_copy {
		err := os.Link(tap, bin)
		if err == nil {
			return 0
		}
	}
	//softlink fallback
	err := os.Symlink(tap, bin)
	if err == nil {
		return 0
	}
	fmt.Println("Can't link tapfile")
	return -2
}

func untap(bin string) {
	//TODO FIX
	//???
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

func walkdotapfn(path string, fi os.FileInfo, err error) error {
	return walkfn(path, fi, err, true)
}
func walkuntapfn(path string, fi os.FileInfo, err error) error {
	return walkfn(path, fi, err, false)
}

func walkfn(path string, fi os.FileInfo, err error, tap bool) error {
	if !fi.IsDir() {
		// if the first character is a ".", then skip it as it's a hidden file

		if strings.HasPrefix(fi.Name(), ".") {
			return nil
		}

		if tap {
			dotap(path)
		} else {
			untap(path)
		}

		//		fmt.Println(path)
		return nil
	}
	return nil
}

func main() {
	do_tapbin("bin")
	filepath.Walk("bin", walkdotapfn)

	//	rm_tapbin("bin")
	//	filepath.Walk("bin", walkuntapfn)

	fmt.Println("Hello world")
}
