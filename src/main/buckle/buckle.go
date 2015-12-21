package main

import (
	"buckle/dirutils"
	"fmt"
	"os/user"
)

func main() {
	fmt.Println("Loading current list...")

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("User's home directory is : ", usr.HomeDir)

	//ioutil.ReadFile(".buckle)

	files, err := dirutils.ListFilesIn("/home/realbot/tmp1")
	if err == nil {
		for i, each := range files {
			fmt.Printf("%d %s\n", i, each)
		}
	} else {
		fmt.Println(err)
	}
}
