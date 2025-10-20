package main

import (
	"fmt"
//	"net/http"
)

const (
	args := os.Args
)

var (
	cat string
	tag string
	help bool
)

func main() {
	fmt.Printf("foo\n")
	for (var i = 0; i < len(args); i++) {
		switch (args[i]) {
		case "-c": //category arg
			cat = args[i+1]
			break
		case "-t": //tag arg
			cat = args[i+1]
			break
		case "-h": //help arg
			help = true
			break
		default:
			//do nothing
		}
}
