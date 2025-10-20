package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "net/http"
    "encoding/json"
    "bytes"
    "github.com/tidwall/gjson"
    "io/ioutil"
)


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
        if v.Get("category").String() == category{
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


/*******************/
/*  main function  */
/*******************/
func main() {
    //get the command arguments
    args := os.Args[1:]
    

    //get the which webhook to use from lastUsedWebhook.txt 
    whichWebhookByteArray, err := os.ReadFile("lastUsedWebhook.txt")
    if err != nil {
        log.Fatalf("err reading lastUsedWebhook.txt:  ", err)
    }
    
    //converts []byte to string then integer
    whichWebhook, err := strconv.Atoi(string(whichWebhookByteArray[0]))
    if err != nil {
        log.Fatalf("err converting []byte to string then integer for the getting which webhook to use", err)
    }
    
    //make sure the webhook used never exceeds the 3rd value
    if whichWebhook == 2 {
        whichWebhook = 0
    } else {
        //switch to next webhook
        whichWebhook++
    }

    //define the data
    category := string(args[0])
    tags := string(args[1])
    text := string(args[2])

    //get the webhook
    webhook := getWebhook(category, int(whichWebhook))

    fmt.Printf("sending timed note:\n")
    fmt.Printf("  args:     %s\n", args)
    fmt.Printf("  webhook:  %s\n", webhook)
    fmt.Printf("  category: %s\n", category)
    fmt.Printf("  tags:     %s\n", tags)
    fmt.Printf("  text:     %s\n", text)
    sendDiscord(text, webhook, tags)
 
    //keep track of current webhook
    os.WriteFile("lastUsedWebhook.txt", []byte(strconv.Itoa(whichWebhook)), 0644)
}


/**************************/
/*  send note to discord  */
/**************************/
func sendDiscord(text string, webhookURL string, noteTag string) {
    //payload struct
    payload := map[string]interface{}{
        "content": "**" + noteTag[:len(noteTag)-2] + " :**\n      " + text + "",
    }

    //convert payload struct to json
    data, err := json.Marshal(payload)
    if err != nil {
        log.Fatalf("err converting payload struct to json:", err)
        os.Exit(1)
    }

    //create the request
    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
    if err != nil {
        fmt.Println("\nsendDiscord()\n  err sending request:\n", err)
        fmt.Println("\n  webhook:\n", webhookURL)
        fmt.Println("\n  text:\n", text)
        fmt.Println("\n  noteTag:\n", noteTag)
        os.Exit(1)
    }

    defer resp.Body.Close()
}
