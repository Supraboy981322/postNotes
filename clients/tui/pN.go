package main   //primary script in postNotes tui client

import (       //import packages
	"fmt"      //  formatted I/O
	"os"       //  os interfacing
	"slices"   //  for slices (clearly)
	"net/http" //  networking over http
	"time"     //  time (clearly)
	"strings"  //  ever-so-slightly more advanced string manipulation
)

const (                            //global const
	version string = "v.0.75.2"    //  for client version
	blue string = "\033[0;34m"     //  for term color blue
	green string = "\033[0;32m"    //  for term color green
	red string = "\033[0;31m"      //  for term color red
	purple string = "\033[0;35m"   //  for term color purple
	yellow string = "\033[0;33m"   //  for term color yellow
	black string = "\033[0;30m"    //  for term color black
	cyan string = "\033[0;36m"     //  for term color cyan
	grey string = "\033[1;30m"     //  for term color grey
	coff string = "\033[0m"        //  to clear term color
)

var (                            //global vars
	cat string                   //  for category arg
	tag string                   //  for tag arg
	data string                  //  for data arg
	help bool                    //  to check if help arg passed
	args []string = os.Args[1:]  //  the args passed
)

func main() {
	var used []int                         //used later to keep track of which args are already checked/used
	var check []bool                       //used later to check if all necessary args are valid
	for i := 0; i < len(args); i++ {       //iterate through all args
		switch (args[i]) {                 //
		case "-c", "--cat", "--category":  //  category arg
			cat = args[i+1]                //    value is next arg (for `-c "foo"`, use `foo` not `-c`)
			used = append(used, i, i+1)    //    note current & next args have been parsed
			break                          //    skip to next iteration
		case "-t", "--tag":                //  tag arg
			tag = args[i+1]                //    value is next arg (for `-t "foo"`, use `foo` not `-t`)
			used = append(used, i, i+1)    //    note current & next args have been parsed
			break                          //    skip to next iteration
		case "-v", "--version":            //  version arg
			fmt.Println(version)           //    print the version number
			used = append(used, i)         //    note current arg has been parsed
			break                          //    skip to next iteration
		case "-h", "--help":               //  help arg
			help = true                    //    note that it's been passed
			used = append(used, i)         //    note current arg has been parsed
			break                          //    skip to next iteration
		default:                           //  if arg doesn't match any of the above... 
			if !slices.Contains(used, i) { //    first, make sure that the arg hasn't been parsed yet
				data = args[i]             //    then, only if it hasn't, use it as the data
			}                              //
		}                                  //
	}                                      //
                                           //
	if help {                              //if the help arg has been passed...
		printHelp()                        //  print the help arg
	}                                      //
	                                       //
	params := []string{cat, tag, data}     //create list of the note's parameters
	for _, arg := range params {           //iterate through all parameters
		if arg != "" {                     //  if the arg exists (will update to a real validity check)...
			check = append(check, true)    //    note that it's valid
		} else {                           //  otherwise...
			check = append(check, false)   //    note that it is invalid
		}                                  //
	}                                      //
	if check[0] && check [1] && check [2] {//if all parameters are valid...
		err := sendNote()                  //  send the note
		if err != nil {                    //  if there was a problem sending the note
			fmt.Println(err)               //    print the problem
		}                                  //
	} else if !help {                      //if any of the parameters are invalid...
		fmt.Println("not enough args")     //  assume not enough args passed (will update to print problem)
		printHelp()                        //  print the help screen
	}
}

func sendNote() error {                                     //send a note
	client := &http.Client{                                 //  create the client instance
		Timeout: time.Second * 10,                          //    timeout after 10 seconds
	}                                                       //
	                                                        //
	request, err := http.NewRequest(                        //  create the request 
		"POST", "[redacted url]",                           //    set to POST and use provided url
		strings.NewReader(data))                            //    use the data param as the body 
	if err != nil {                                         //    if there is a problem...
		return fmt.Errorf("err creating http request", err) //      return it
	}                                                       //
                                                            //
	request.Header.Add("Content-Type", "text/plain")        //  set the `Content-Type` header
	request.Header.Add("c", cat)                            //  set the category header
	request.Header.Add("t", tag)                            //  set the tag header
                                                            //
	response, err := client.Do(request)                     //  send the request
	if err != nil {                                         //    if there is a problem...
		return fmt.Errorf("err sending http request", err)  //      return it
	}                                                       //
	defer response.Body.Close()                             //  close the body
	return nil                                              //  presumed to not have had any problems
}

func printHelp() {                                                  //```help
	fmt.Printf("%susage:%s\n", cyan, coff)                          //usage:\n
	fmt.Printf("%s  %scategory:%s\n", grey, yellow, coff)           //  category:\n
	fmt.Printf("%s    %s`%s-c%s`,", grey, coff, yellow, coff)       //    `-c`,
	fmt.Printf("%s %s`%s--cat%s`,", grey, coff, yellow, coff)       //    `--cat`,
	fmt.Printf("%s %s`%s--category%s`\n", grey, coff, yellow, coff) //    `--category`\n
	fmt.Printf("%s  %stag:%s\n", grey, green, coff)                 //  tag:\n
	fmt.Printf("%s    %s`%s-t%s`,", grey, coff, green, coff)        //    `-t`,
	fmt.Printf("%s %s`%s--tag%s`\n", grey, coff, green, coff)       //    `--tag`\n
	fmt.Printf("%s  %shelp:%s\n", grey, purple, coff)               //  help:\n
	fmt.Printf("%s    %s`%s-h%s`,", grey, coff, purple, coff)       //    `-h`,
	fmt.Printf("%s  %s`%s--help%s`\n", grey, coff, purple, coff)    //    `--help`\n
	fmt.Printf("%s  %sversion:%s\n", grey, red, coff)               //  version:
	fmt.Printf("%s    %s`%s-v%s`,", grey, coff, red, coff)          //    `-v`,
	fmt.Printf("%s  %s`%s--version%s`\n", grey, coff, red, coff)    //    `--version`\n
}                                                                   //```
