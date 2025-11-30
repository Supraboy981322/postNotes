package main

import (
//	"net/http"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	port int
	conf gomn.Map
	cats map[string]map[string][]string
	lastUsedWebhook = 0;
	logLvl = log.DebugLevel
)

func init() {
	log.Info("starting...")
	log.SetLevel(logLvl)

	log.Info("configuring...")
	config()
	log.Info("configured")
}

func main() {
	
}
