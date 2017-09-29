package main

import (
	"fmt"
	. "github.com/mattn/go-getopt"
	"github.com/mudler/devbot/bot"
	"log"
	"os"
	"time"

	_ "github.com/mudler/devbot/plugins/admin"
	_ "github.com/mudler/devbot/plugins/brain"
	_ "github.com/mudler/devbot/plugins/ddg"
	_ "github.com/mudler/devbot/plugins/helper"

	_ "github.com/mudler/devbot/plugins/gentoobug"
	_ "github.com/mudler/devbot/plugins/logger"
	_ "github.com/mudler/devbot/plugins/progresstracker"
	_ "github.com/mudler/devbot/plugins/rss"
	_ "github.com/mudler/devbot/plugins/sabayonbug"
	_ "github.com/mudler/devbot/plugins/seen"
	_ "github.com/mudler/devbot/plugins/spamdetect"
	_ "github.com/mudler/devbot/plugins/team"

	_ "github.com/mudler/devbot/plugins/urlpreview"
)

func main() {

	var c int
	var configurationFile = "default.json"
	var logFile string
	OptErr = 0
	for {
		if c = Getopt("c:l:h"); c == EOF {
			break
		}
		switch c {
		case 'c':
			configurationFile = OptArg
		case 'l':
			logFile = OptArg
		case 'h':
			println("usage: " + os.Args[0] + " [-c configfile.json|-l logfile|-h]")
			os.Exit(1)
		}
	}
	fmt.Println("Configuration file: " + configurationFile)
	config, err := bot.LoadConfig(configurationFile)
	if err != nil {
		log.Fatal("Error loading configuration file: " + err.Error())
	}

	if logFile != "" {
		//Set logging to file
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error opening file: %v", err)
		}
		defer f.Close()
		fmt.Println("Log file: " + logFile)

		log.SetOutput(f)
	}
	for {
		bot.Start(config)
		time.Sleep(1000 * time.Millisecond)
	}
}
