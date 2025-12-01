package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"errors"
	"slices"
	"strings"
	"net/http"
//	"path/filepath"
	"github.com/Supraboy981322/gomn"
)

var (
	url string
	cat string
	tag string
	data string
	timeout int
	args = os.Args[1:]
)

func init() {
	var ok bool
	var err error
	var conf gomn.Map

	if conf, err = gomn.ParseFile("conf.gomn"); err != nil {
		erorF("reading config", err)
	}

	var server gomn.Map
	if server, ok = conf["server"].(gomn.Map); !ok {
		erorF("getting server config", err)
	}

	if url, ok = server["address"].(string); !ok {
		err := errors.New("invalid server address")
		erorF("parsing server config", err)
	}

	if timeout, ok = conf["timeout"].(int); !ok {
		err := errors.New("invalid timeout")
		erorF("parsing server config", err)
	}

	var taken []int
	for i, arg := range args {
		if !slices.Contains(taken, i) {
			switch arg {
			case "-t", "-tag", "--tag", "tag":
				tags := strings.Split(args[i+1], ";")
				tag = strings.Join(tags, "`**, **`")
				tag = "**`"+tag+"`**"
				taken = append(taken, i+1)
			case "-c", "-cat", "--cat", "--category", "cat", "category":
				cat = args[i+1]
				taken = append(taken, i+1)
			case "-h", "--help", "help", "-help", "--h":
				help()
			default:
				if data == "" {
					data = arg
				} else if url == "" {
					url = arg
				} else {
					err := errors.New("invalid arg")
					erorF("parsing args", err)
				}
			}
		}
	};for _, thing := range []string{url, cat, tag} {
		if thing == "" {
			err := errors.New("missing arg")
			erorF("parsing args", err)
		}
	}
}

func main() {
	client := &http.Client{
		Timeout: time.Second*time.Duration(timeout),
	}

	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s/post", url),
		strings.NewReader(data))
	if err != nil {
		erorF("creating request", err)
	}

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("c", cat)
	req.Header.Add("t", tag)

	resp, err := client.Do(req)
	if err != nil {
		erorF("sending request", err)
	}; defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		erorF("reading response", err)
	}

	fmt.Println(string(respBody))
}
