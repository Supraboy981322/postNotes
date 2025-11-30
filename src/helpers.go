package main

import (
	"strings"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

func config() {
	var err error
//	var ok bool

	log.Debug("reading config")
	if conf, err = gomn.ParseFile("conf.gomn"); err != nil {
		log.Fatalf("err reading config:  %v", err)
	} else { log.Debug("read config") }

	log.Debug("setting log level")
	if logLvlStrRaw, ok := conf["log level"]; ok {
		if logLvlStr, ok := logLvlStrRaw.(string); ok {
			switch strings.ToLower(logLvlStr)	{
				case "debug":	logLvl = log.DebugLevel
				case "warn":	logLvl = log.WarnLevel
				case "info":  logLvl = log.InfoLevel
				case "error":	logLvl = log.ErrorLevel
				case "fatal":	logLvl = log.FatalLevel
			}
			log.SetLevel(logLvl)
		} else { log.Fatal("err parsing log level") }
	} else { log.Warn("log level not set, defaulting to debug") }

	log.Debug("setting port")
	if portRaw, ok := conf["port"]; ok {
		if port, ok = portRaw.(int); !ok {
			log.Fatal("err setting port from config")
		} else { log.Debug("set port") }
	} else { log.Fatal("err getting port from config") }

	log.Debug("reading categories")
	if catsRaw, ok := conf["categories"]; !ok {
		log.Fatal("err reading categories")
	} else {
		log.Debug("read categories")
		log.Debug("mapping categories")
		catsMapBuilt := make(map[string]map[string][]string)
		if catsMap, ok := catsRaw.(gomn.Map); ok {
			for cat, stuffRaw := range catsMap {
				var aliases, webhooks []string
				var catName string

				if catName, ok = cat.(string); !ok {
					log.Fatalf("err on category name:  %#v", cat)
				} else { log.Debugf("got category name: %s", catName) }

				var stuff gomn.Map
				if stuff, ok = stuffRaw.(gomn.Map); !ok {
					log.Fatalf("err asserting category (\"%s\") config", catName)
				} else { log.Debugf("got category (\"%s\") map", catName) }
				
				var aliasesRaw []interface{}
				if aliasesRaw, ok = stuff["aliases"].([]interface{}); !ok {
					log.Fatal("err asserting aliasesRaw to interface slice")
				} else { log.Debug("asserted aliasesRaw to interface slice") }

				for i, aliasRaw := range aliasesRaw {
					if alias, ok := aliasRaw.(string); !ok {
						log.Fatalf("err asserting category (\"%s\") alias (#%d)", catName, i)
					} else { aliases = append(aliases, alias) }
				}; log.Debug("got category aliases")
				
				var webhooksRaw []interface{}
				if webhooksRaw, ok = stuff["webhooks"].([]interface{}); !ok {
					log.Fatal("err asserting webhooksRaw to interface slice")
				} else { log.Debug("asserted webhooksRaw to interface slice") }

				for i, webhookRaw := range webhooksRaw {
					if webhook, ok := webhookRaw.(string); !ok {
						log.Fatalf("err asserting category (\"%s\") webhooks (#%d)", catName, i)
					} else { webhooks = append(webhooks, webhook) }
				}; log.Debug("got category webhooks")


				catMap := make(map[string][]string)
				catMap["aliases"] = aliases
				catMap["webhooks"] = webhooks
				catsMapBuilt[catName] = catMap

			}; cats = catsMapBuilt
		} else { log.Fatal("failed to parse category map") }
	}
}
