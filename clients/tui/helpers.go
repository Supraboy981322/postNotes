package main

import (
	"os"
	"fmt"
)

func eror(str string, err error) {
	fmt.Printf("err %s:  %v\n", str, err)
}

func erorF(str string, err error) {
	eror(str, err)
	os.Exit(1)
}

func help() {
	lines := []string{
		"postNotes --> help",
		"  usage",
		"    help",
		"      \033[0;36m-h\033[0m, \033[0;36m--help\033[0m",
		"    category",
		"      \033[0;31m-c\033[0m, \033[0;31m--cat\033[0m",
		"    tag",
		"      \033[0;32m-t\033[0m, \033[0;32m--tag\033[0m",
		"  example",
		"    pN \033[0;32m-c \"todo\"\033[0m"+
					"\033[0;31m-t \"NixOS\"\033[0m"+
					" \"fix app icons on desktop\"",
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
