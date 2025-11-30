package main

import (
	"fmt"
	"io"
    "io/ioutil"
	"log"
	"net/http"
	"os"
    //"os/exec"
	"bytes"
	"encoding/json"
    "sort"
    "github.com/tidwall/gjson"
    "strings"
    "strconv"
//    "time"
//    "github.com/go-co-op/gocron/v2"
)

type (
	noteStruct struct {
    NOTEtag []string `json:"tag"`
    NOTEcategory string `json:"category"`
    NOTEtext string `json:"text"`
	}
	categoryStruct struct {
    CATEGORYname string `json:"category"`
    CATEGORYshorthands []string `json:"shorthands"`
    CATEGORYwebhooks []string `json:"webhook"`
	}
)

var (
	categories []categoryStruct
)
//var scheduler gocron.NewScheduler


/**********************************/
/*  handle setup of a timed note  */
/**********************************/
/*func timeNote(text string, tags string, category string, noteTimeSlice []string) {
    if len(noteTimeSlice) < 4 {
        log.Fatalf("err note time slice length is less than 4, THIS IS A CRITICAL ERROR WHICH SHOULD NEVER OCCUR:  %#v", noteTimeSlice)
    }

    noteTimeSliceClock := strings.Split(noteTimeSlice[0], "-")
   
    if len(noteTimeSliceClock) < 2 {
        log.Fatalf("invalid clock format: %q", noteTimeSlice[0])
    }

    minute := noteTimeSliceClock[1]
    hour := noteTimeSliceClock[0]
    day := noteTimeSlice[1]
    month := noteTimeSlice[2]
    dayOfWeek := noteTimeSlice[3]

    fmt.Printf("min:  %s\nhour:  %s\nday:  %s\nmonth:  %s\ndayOfWeek%s\n", minute, hour, day, month, dayOfWeek)

    //create time string
    noteTimeString := fmt.Sprintf("%s %s %s %s %s", minute, hour, day, month, dayOfWeek)
    fmt.Printf("noteTimeString = %s\n", noteTimeString)

    //add cronjob to scheduler
    //_, err := scheduler.NewJob(
    //    gocron.CronJob(noteTimeString, false),
    //    gocron.NewTask(func() {
    //    scheduler.Cron(noteTimeString).Do(sendFunc(text, category, tags)),
    //)
    //if err != nil {
    //    log.Fatalf("err attempting add postNote to scheduler:  %v", err)
    //}    
}*/



func sendFunc(text string, category string, tags string) {
    //for debugging
    log.Println("if you see this, the function inside the timed note gocron is executing")

    //create temp. var for webhook
    var webHookURL string
            
    //get the which webhook to use from lastUsedWebhook.txt
    whichWebhookByteArray, err := os.ReadFile("lastUsedWebhook.txt")
    if err != nil {
        log.Fatalf("err reading lastUsedWebhook.txt:  ", err)
    }
            
    //converts []byte to string then integer
    whichWebhook, err := strconv.Atoi(string(whichWebhookByteArray[0]))
    if err != nil {
        log.Fatalf("err converting []byte to string then integer for the getting which webhook to use in timed note:  ", err)
    }
        
    //make sure that the value for the webhook used doesn't exceed the 3rd webhook
    if whichWebhook == 2 {
        whichWebhook = 0
    } else {//switch to next webhook
        whichWebhook++
    }

    //keep track of current webhook
    os.WriteFile("lastUsedWebhook.txt", []byte(strconv.Itoa(whichWebhook)), 0644)

    //get the webhook from categories.json
    webHookURL = getWebhook(category, whichWebhook)

    //send it to discord
    sendDiscord(text, webHookURL, tags)
}


/*****************************/
/*  handle webpage requests  */
/*****************************/
func webPageHandler(w http.ResponseWriter, r *http.Request) {
    //get the requested page
    requestedPage := r.URL.Path

    //serve only these webpages
    var webpageFileName string
    if requestedPage == "/script.js" {
        w.Header().Set("Content-Type", "text/javascript")
        webpageFileName = "clients/web/script.js"
    } else if requestedPage == "/notes.json" {
        webpageFileName = "notes.json"        
    } else if requestedPage == "/main.css" {
        w.Header().Set("Content-Type", "text/css")
        webpageFileName = "clients/web/main.css"
    } else {
        webpageFileName = "clients/web/index.html"
    }
    
    //read the file's contents
    webpageContent, err := ioutil.ReadFile(webpageFileName)
    if err != nil {
        log.Fatalf("err reading %s:  %s", webpageFileName, err)
    }

    //return the file's contents
    fmt.Fprintf(w, "%s", string(webpageContent))
}


/********************************/
/*  handle creating categories  */
/********************************/
func createCatHandler(w http.ResponseWriter, r *http.Request) {
    //let user know that the category is being created
    fmt.Printf("attempting to create category")
    w.Write([]byte("category creation request recieved\n"))

    //define the category json file name
    categoriesJSONfileName := "categories.json"

    //read categories.json for current data
    currentCategoriesJSON, err := ioutil.ReadFile(categoriesJSONfileName)
    if err != nil {
        log.Fatalf("err reading categories.json file for category creation:  ", err)
    }

    //get the category's data
    categoryName := r.Header.Get("name")
    categoryShorthandsRaw := r.Header.Get("shorthand")
    categoryWebhooksRaw := r.Header.Get("webhooks")
    //split the category shorthands into a slice
    categoryShorthands := strings.Split(categoryShorthandsRaw, ";")
    //split the category's webhooks into a slice
    categoryWebhooks := strings.Split(categoryWebhooksRaw, ";")
    //add the full category name to the shorthands list
    categoryShorthands = append(categoryShorthands, categoryName)

    //unmarshal json for current categories data
    var categoriesSlice []categoryStruct
    err = json.Unmarshal(currentCategoriesJSON, &categoriesSlice)
    if err != nil {
        log.Fatalf("err unmarshalling categories json data for categories creation:  ", err)
    }

    //define the new category's struct
    categoryData := categoryStruct{
        CATEGORYname: categoryName,
        CATEGORYshorthands: categoryShorthands,
        CATEGORYwebhooks: categoryWebhooks,
    }

    //add new category to categories slice
    categoriesSlice = append(categoriesSlice, categoryData)

    //marshal updated categories slice to new JSON
    newCategoriesJSON, err := json.MarshalIndent(categoriesSlice, "", " ")
    if err != nil {
        log.Fatalf("err marshalling categoryData for creating category:  ", err)
    }

    //update the categories.json file
    err = ioutil.WriteFile(categoriesJSONfileName, newCategoriesJSON, 0644)
    if err != nil {
        log.Fatalf("err attempting to write updated category data to categories.json for category creation:  ", err)
    }

    fmt.Printf("wrote new category to json file")
    w.Write([]byte("category creation successful"))
    w.Write([]byte("  name:  \n    " + categoryName + "\n"))
    w.Write([]byte("  shorthands: \n    " + strings.Join(categoryShorthands, "\n    ") + "\n"))
    w.Write([]byte("  webhooks:  \n    " + strings.Join(categoryWebhooks, "\n    ") + "\n"))
}

/***************************/
/*  handle fetching notes  */
/***************************/
func getNoteHandler(w http.ResponseWriter, r *http.Request) {
    //define the name of the notes' json file
    notesFileName := "notes.json"
    
    //read the contents of the notes' json file
    currentNotes, err := ioutil.ReadFile(notesFileName)
    if err != nil {
        log.Fatalf("err reading notes' json file:  ", err)
    }

    //turn the notes into a struct by unmarshalling it
    var notes []noteStruct
    err = json.Unmarshal(currentNotes, &notes)
    if err != nil {
        log.Fatalf("err reading note struct")
    }

    //sort the notes by category, then sort them by tag
    sort.Slice(notes, func(i, j int) bool {
        if notes[i].NOTEcategory == notes[j].NOTEcategory{
            return strings.Join(notes[i].NOTEtag, ",") < strings.Join(notes[j].NOTEtag, ",")
        }
        return notes[i].NOTEcategory < notes[j].NOTEcategory
    })
	for _, note := range notes {
    fmt.Fprintf(w, "category:  %s\n", note.NOTEcategory)
    fmt.Fprintf(w, "..tag:       %s\n", note.NOTEtag)
    fmt.Fprintf(w, "....content:   %s\n", note.NOTEtext)
  	fmt.Fprintf(w, "\n")
  }
}


/****************************/
/*  get webhooks from json  */
/****************************/
func getWebhook(category string, which int) string {
  //get the webhook categories.json file contents
  categoriesJSONbyte, err := ioutil.ReadFile("categories.json")
  if err != nil {
  	log.Fatalf("err reading categories.json:  ", err)
  }

  //convert byte[] to string
  categoriesJSON := string(categoriesJSONbyte)
    
  index := -1
  gjson.Parse(categoriesJSON).ForEach(func(i, v gjson.Result) bool {
    if v.Get("category").String() == category {
      //get the index
      index = int(i.Int())

      //exit current iteration of loop
    	return false
    }
    //move to next iteration of loop
    return true
  })

  if index != -1 {
    webhookPath := fmt.Sprintf("%d.webhook.%d", index, which)

    fmt.Printf("webhookPath: %s\n", webhookPath)
    webhook := gjson.Get(categoriesJSON, webhookPath).String()
        
    //return the webhook
  	return webhook
  } else {
  	return " "
  }
}


/***********************************/
/*  converts shorthand categories  */
/*      to full category names     */
/***********************************/
func convertShortHand(shorthand string) string {
  //define the categories values
  categoriesJSONfileName := "categories.json"

  //read the categories json file
  categoriesJSONbyte, err := ioutil.ReadFile(categoriesJSONfileName)
  if err != nil {
  	log.Fatalf("err reading categories json file:  ", err)
  }

  //convert the []byte json data to string
  categoriesJSON := string(categoriesJSONbyte)

  //set the category to blank
  category := ""
  gjson.Parse(categoriesJSON).ForEach(func(_, v gjson.Result) bool {
    found := false
    v.Get("shorthands").ForEach(func(_, s gjson.Result) bool {
      if s.String() == shorthand {
	    	//set the index where it was found
	      category = v.Get("category").String()
	      //set found
	      found = true
	      //stop shorthand loop
		    return false
		  }
		  return true
		})
	  if found {
			//stop loop
	  	return false
		}
    return true
  })
  //return the index of the shorthand
  return category
}


/***********************/
/* save notes locally  */
/***********************/
func saveNote(text string, category string, tags []string) {
  //define the name of the notes json file
  notesFileName := "notes.json"
    
  //read the current notes json file
  currentNotes, err := ioutil.ReadFile(notesFileName)
  if err != nil {
  	log.Fatalf("err reading notes' json file:  ", err)
  }

  //create an array for the notes
  var notesArray []noteStruct
  //unmarshal the notes' json and put it into the array
  err = json.Unmarshal(currentNotes, &notesArray)
  if err != nil {
    log.Fatalf("err attempting to unmarshal notes array", err)
  }

  //define the new note
  newNote := noteStruct{
    NOTEtag: tags, //tags slice from input
    NOTEcategory: category,
  	NOTEtext: text,
  }

  //add the new note to the notes' array
  notesArray = append(notesArray, newNote)
    
  //marshal the notes array to json
  updatedNotesJSON, err := json.MarshalIndent(notesArray, "", "  ")
  if err != nil {
  	log.Fatalf("err attempting to marshal notes array back into JSON:  ", err)
  }
    
  //write the updated notes' json to the notes' json file
  err = ioutil.WriteFile(notesFileName, updatedNotesJSON, 0644)
  if err != nil {
  	log.Fatalf("err writing the notes' json to file", err)
  }

  fmt.Printf("wrote new note to json file\n")
}


/***************************/
/*  note creation handler  */
/***************************/
func noteHandler(w http.ResponseWriter, r *http.Request) {
	//protocolScheme := r.Header.Get("X-Forwarded-Proto")
	if r.Method != http.MethodPost {
		fmt.Printf("only post allowed")
    w.Write([]byte("only post allowed\n"))
		return
	}
	// Read body
	body, err := io.ReadAll(r.Body)

	//if error, return error
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	//close the body
	defer r.Body.Close()

	//convert request body to string
	text := string(body)

	//confirm receival of note
  w.Write([]byte("note recieved:  " + text + "\n"))
	
	//create temp var for category
	var category string
	var noteTag string
  var timed bool
  if r.Header.Get("time") != "" {
    timed = true
  	w.Write([]byte("running as timed note"))
  } else {
  	timed = false
  }

	w.Write([]byte("  checking category:\n"))
	//get the category
	if r.Header.Get("c") != "" {
		category = convertShortHand(r.Header.Get("c"))
	} else if r.Header.Get("cat") != "" {
		category = convertShortHand(r.Header.Get("cat"))
	} else if r.Header.Get("category") != "" {
		category = convertShortHand(r.Header.Get("category"))
	} else {
		//if no category specified...
		category = "misc"
        
    //respond...

		//make sure user knows that means it defaults to `misc`
		w.Write([]byte("  no category specified; defaulting to `misc`\n"))
        
		//make sure the user knows how to specify a category
		w.Write([]byte("    valid headers to specify category:\n"))
		w.Write([]byte("      c: [your category]\n"))
		w.Write([]byte("      cat: [your category]\n"))
		w.Write([]byte("      category: [your category]\n"))
	}

	//get the tag
    var noteTagsRaw string
	if r.Header.Get("t") != "" {
		noteTagsRaw = r.Header.Get("t")
	} else if r.Header.Get("tag") != "" {
		noteTagsRaw = r.Header.Get("tag")
	} else if r.Header.Get("notetag") != "" {
		noteTagsRaw= r.Header.Get("notetag")
	} else {
		//if no tag specified...
		noteTagsRaw = "misc"

		//make sure user knows that means it defaults to `misc`
		w.Write([]byte("  no note tag specified; not using a note tag\n"))

		//make sure the user knows how to specify a noteTag
		w.Write([]byte("    valid headers to specify tag:\n"))
		w.Write([]byte("      t: [your tag]\n"))
		w.Write([]byte("      tag: [your tag]\n"))
		w.Write([]byte("      noteTag: [your tag]"))
	}

  //convert the note tags string to a slice
  noteTagsSlice := strings.Split(noteTagsRaw, ";")

  var noteTagsFormatted string
  for _, value := range noteTagsSlice {
  	noteTagsFormatted += "`" + value + "` , "
  }

	if timed {
		w.Write([]byte("TODO: timed notes"))
  } else {
		//create temp var for the webhook url
    var webHookURL string
        
    //get the which webhook to use from lastUsedWebhook.txt 
    whichWebhookByteArray, err := os.ReadFile("lastUsedWebhook.txt")
    if err != nil {
      log.Fatalf("err reading lastUsedWebhook.txt:  ", err)
    }
        
    fmt.Printf(string(whichWebhookByteArray))

    //converts []byte to string then integer
    whichWebhook, err := strconv.Atoi(string(whichWebhookByteArray[0]))
    if err != nil {
    	log.Fatalf("err converting []byte to integer for webhook:  ", err)
    }
        
    //make sure that the value for the webhook used doesn't exceed the 3rd webhook
    if whichWebhook == 2 {
    	whichWebhook = 0
    } else {
    	//switch to next webhook
      whichWebhook++
    }

    //keep track of current webhook
    os.WriteFile("lastUsedWebhook.txt", []byte(strconv.Itoa(whichWebhook)), 0644)

    //get the webhook from categories.json
    webHookURL = getWebhook(category, whichWebhook)

    if webHookURL == " " {
    	fmt.Printf("Invalid category")
      w.Write([]byte("Invalid category"))
    } else {  
			//respond with data
			w.Write([]byte("    category:  " + category + "\n"))
			w.Write([]byte("    tag:       " + noteTag + "\n"))

			//print header and body to console
			fmt.Printf("  category:  %s\n", category)
			fmt.Printf("  tag:       %s\n", noteTagsFormatted)
			fmt.Printf("  content:   %s\n", text)
      
      //save note to json file
      saveNote(text, category, noteTagsSlice)
      fmt.Printf("webhook: %s\n", webHookURL)

			//send it to discod
			sendDiscord(text, webHookURL, noteTagsFormatted)

			//response
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("noted\n"))
    }
  }
}

/*func init() {
  var err error
  scheduler, err = gocron.NewScheduler()
  if err != nil {
    log.Fatalf("failed to create scheduler: %v", err)
  }
  scheduler.Start()

  //defer scheduler.Shutdown()
}*/

/*******************/
/*  main function  */
/*******************/
func main() {
/*  var err error
  scheduler, err = gocron.NewScheduler()
  if err != nil {
    log.Fatalf("failed to create scheduler: %v", err)
  }
  scheduler.Start()*/
	//handle post requests for notes
	http.HandleFunc("/post", noteHandler)

  //handle get requests for notes
  http.HandleFunc("/get", getNoteHandler)

  //handle creating categories
  http.HandleFunc("/createCategory", createCatHandler)

  //handle all other pages
  http.HandleFunc("/", webPageHandler)

  //specify the port to listen on
	port := "6502"

  //print that it's working to console
	log.Printf("Listening on http://localhost:%s (POST requests)\n", port)

  if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}

//    defer scheduler.Shutdown()
}

/**************************/
/*  send note to discord  */
/**************************/
func sendDiscord(text string, webHookURL string, noteTag string) {
	//payload struct
	payload := map[string]interface{}{
    "content": "**" + noteTag[:len(noteTag)-2] + " :**\n      " + text + "",
	}

	//convert payload struct to json
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error converting payload struct to json:", err)
		os.Exit(1)
	}

	//create the request
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
    fmt.Println("\nsendDiscord()\n  error sending request:\n", err)
    fmt.Println("\n  webhook:\n", webHookURL)
    fmt.Println("\n  text:\n", text)
    fmt.Println("\n  noteTag:\n", noteTag)
		os.Exit(1)
	}
	defer resp.Body.Close()
}

