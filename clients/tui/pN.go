package main

import (
	"fmt"
	"os"
	"slices"
	"net/http"
	"time"
	"strings"
)

const (
	version string = "v.0.75.0"
	blue string = "\033[0;34m"
	green string = "\033[0;32m"
	red string = "\033[0;31m"
	purple string = "\033[0;35m"
	yellow string = "\033[0;33m"
	black string = "\033[0;30m"
	cyan string = "\033[0;36m"
	grey string = "\033[1;30m"
	coff string = "\033[0m"
)

var (
	check []bool
	cat string
	tag string
	data string
	help bool
	args []string = os.Args[1:]
)

func main() {
	var used []int
	for i := 0; i < len(args); i++ {
		switch (args[i]) {
		case "-c", "--cat", "--category": //category arg
			cat = args[i+1]
			used = append(used, i, i+1)
			break
		case "-t", "--tag": //tag arg
			tag = args[i+1]
			used = append(used, i, i+1)
			break
		case "-v", "--version":
			fmt.Println(version)
			used = append(used, i)
			break
		case "-h", "--help": //help arg
			help = true
			used = append(used, i)
			break
		default:
			if !slices.Contains(used, i) {
				data = args[i]
			}
		}
	}

	if help {
		printHelp()
	}

	for _, arg := range []string{cat, tag, data} {
		if arg != "" {
			check = append(check, true)
		} else {
			check = append(check, false)
		}
	}
	if check[0] && check [1] && check [2] {
		if err := sendNote(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("not enough args")
		printHelp()
	}
}

func sendNote() error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	
	request, err := http.NewRequest(
		"POST", "[redacted url]",
		strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("err creating http request", err)
	}

	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("c", cat)
	request.Header.Add("t", tag)

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("err sending http request", err)
	}
	defer response.Body.Close()
	
	fmt.Println("response status", response.Status)
	return nil
}

func printHelp() {
	fmt.Printf("%susage:%s\n", cyan, coff)
	fmt.Printf("%s  %scategory:%s\n", grey, yellow, coff)
	fmt.Printf("%s    %s`%s-c%s`,", grey, coff, yellow, coff)
	fmt.Printf("%s %s`%s--cat%s`,", grey, coff, yellow, coff)
	fmt.Printf("%s %s`%s--category%s`\n", grey, coff, yellow, coff)
	fmt.Printf("%s  %stag:%s\n", grey, green, coff)
	fmt.Printf("%s    %s`%s-t%s`,", grey, coff, green, coff)
	fmt.Printf("%s %s`%s--tag%s`\n", grey, coff, green, coff)
	fmt.Printf("%s  %shelp:%s\n", grey, purple, coff)
	fmt.Printf("%s    %s`%s-h%s`,", grey, coff, purple, coff)
	fmt.Printf("%s  %s`%s--help%s`\n", grey, coff, purple, coff)
	fmt.Printf("%s  %sversion:%s\n", grey, red, coff)
	fmt.Printf("%s    %s`%s-v%s`,", grey, coff, red, coff)
	fmt.Printf("%s  %s`%s--version%s`\n", grey, coff, red, coff)
}
