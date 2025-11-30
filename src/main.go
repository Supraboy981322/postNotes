package main

import (
	"io"
	"bytes"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	port int
	conf gomn.Map
	logLvl = log.DebugLevel
	lastUsedWebhooks map[string]int
	cats map[string]map[string][]string
)

func init() {
	log.Info("starting...")
	log.SetLevel(logLvl)

	log.Info("configuring...")
	config()
	log.Info("configured")
}

func main() {
	http.HandleFunc("/post", noteHandler)
//	http.HandleFunc("/get", getNotesHandler)
	http.HandleFunc("/createCategory", createCatHandler)
//	http.HandleFunc("/", webPageHandler)

	log.Infof("listening on port:  %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func createCatHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("creating category")
	w.Write([]byte("creation request recieved..."))

	w.Write([]byte("TODO:  category creation"))
	log.Warn("TODO:  category creation")

	w.Write([]byte("request processed"))
}

func getWebhook(cat string) string {
	which := lastUsedWebhooks[cat]+1
	if which > 2 {
		which = 0
	}
	
	lastUsedWebhooks[cat] = which

	return cats[cat]["webhooks"][which]
}

func convertAlias(shorthand string) string {
	var fullname string
	var found bool

	for cat, stuff := range cats {
		if !found {
			aliases := stuff["aliases"]
			for _, alias := range aliases {
				if alias == shorthand || shorthand == cat {
					found = true
					fullname = cat
				}
			}
		} else { break }
	}

	return fullname
}

func saveNote(txt string, cat string, tags []string) {
	log.Warn("TODO:  saving notes")
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("only POST method allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	};defer r.Body.Close()

	txt := string(body)
	w.Write([]byte("note recieved:  "+txt+"\n"))

	var cat string
//	var timed bool
	if r.Header.Get("time") != "" {
	//	timed = true
		w.Write([]byte("TODO:  timed notes"))
	}

	for _, chk := range []string{"c", "cat", "category"} {
		catChk := r.Header.Get(chk)
		if catChk != "" && cat == "" {
			cat = convertAlias(catChk)
		}
	}; if cat == "" {
		w.Write([]byte("invalid category, defaulting to \"misc\""))
		cat = "misc"
	}
	
	var tagsRaw string
	for _, chk := range []string{"t", "tag", "notetag", "tags"} {
		tagChk := r.Header.Get(chk)
		if tagChk != "" && tagsRaw == "" {
			tagsRaw = tagChk
		}
	}; if tagsRaw == "" {
		w.Write([]byte("invalid category, defaulting to \"misc\""))
		tagsRaw = "misc"
	}

	tagsSlice := strings.Split(tagsRaw, ";")

	var tagsFmt string
	for _, tag := range tagsSlice {
		tagsFmt += "`" + tag + "` , "
	}
	
	webhook := getWebhook(cat)

	log.Debugf("cat:  %s", cat)
	sendDiscord(txt, webhook, tagsFmt)
}

func sendDiscord(txt string, webhook string, tags string) {
	payload := map[string]interface{}{
		"content": "**" + tags[:len(tags)-2] + " :**\n      " + txt,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("marshalling payload to json:  %v", err)
	}

	buf := bytes.NewBuffer(data)
	cType := "application/json"
	resp, err := http.Post(webhook, cType,	buf)
	if err != nil {
	log.Errorf("sending to discord:  %v", err)
	};defer resp.Body.Close()
}
