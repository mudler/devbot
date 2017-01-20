package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	. "github.com/mattn/go-getopt"
	"github.com/mudler/devbot/shared/registry"
	"github.com/mudler/devbot/shared/utils"
	"github.com/thoj/go-ircevent"

	_ "github.com/mudler/devbot/plugins/admin"
	_ "github.com/mudler/devbot/plugins/brain"
	_ "github.com/mudler/devbot/plugins/ddg"
	_ "github.com/mudler/devbot/plugins/helper"

	_ "github.com/mudler/devbot/plugins/gentoobug"
	_ "github.com/mudler/devbot/plugins/logger"
	_ "github.com/mudler/devbot/plugins/sabayonbug"
	_ "github.com/mudler/devbot/plugins/spamdetect"

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
	config, err := util.LoadConfig(configurationFile)

	if logFile != "" {
		//Set logging to file
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	conn := irc.IRC(config.BotNick, config.BotUser)
	conn.UseTLS = config.Tls

	fmt.Println("BotNick:\t" + config.BotNick)
	fmt.Println("BotUser:\t" + config.BotUser)
	fmt.Println("Channels:")
	for _, channel := range config.Channel {
		fmt.Println("\t" + channel)
	}

	if err != nil {
		fmt.Println("Failed to connect.")
		panic(err)
	}

	plugin_registry.Config = config
	plugin_registry.Conn = conn

	// Bootstrapper for plugins
	for _, d := range plugin_registry.Plugins {
		if (len(config.Plugins) > 0 && config.EnabledPlugin(plugin_registry.KeyOf(d))) || len(config.Plugins) == 0 {
			go d.Register()
		}
	}

	// Bootstrapper for plugins
	for _, d := range plugin_registry.Plugins {
		if len(config.Plugins) > 0 && !config.EnabledPlugin(plugin_registry.KeyOf(d)) {
			plugin_registry.DisablePlugin(plugin_registry.KeyOf(d))
		}
	}

	log.Println(strconv.Itoa(len(plugin_registry.Plugins)) + " plugins loaded")

	if logFile != "" {
		fmt.Println("Log file: " + logFile)
	}

	conn.AddCallback("001", func(e *irc.Event) {
		for _, channel := range config.Channel {
			conn.Join(channel)
		}
		if config.BotNickPassword != "" {
			log.Println("Identifying bot nickname against NickServ")
			conn.Privmsg("NickServ", "IDENTIFY "+config.BotNickPassword)
		}
	})

	conn.AddCallback("PRIVMSG", func(event *irc.Event) {
		for _, d := range plugin_registry.Plugins {
			if obj, ok := d.(interface {
				OnPrivmsg(*irc.Event)
			}); ok {
				go obj.OnPrivmsg(event)
			}
		}
	})
	conn.AddCallback("JOIN", func(event *irc.Event) {
		for _, d := range plugin_registry.Plugins {
			if obj, ok := d.(interface {
				OnJoin(*irc.Event)
			}); ok {
				go obj.OnJoin(event)
			}
		}

	})
	conn.AddCallback("PART", func(event *irc.Event) {
		for _, d := range plugin_registry.Plugins {
			if obj, ok := d.(interface {
				OnPart(*irc.Event)
			}); ok { //Check if plugin has that method
				go obj.OnPart(event)
			}
		}
	})
	conn.AddCallback("QUIT", func(event *irc.Event) {
		for _, d := range plugin_registry.Plugins {
			if obj, ok := d.(interface {
				OnQuit(*irc.Event)
			}); ok { //Check if plugin has that method
				go obj.OnQuit(event)
			}
		}

	})

	err = conn.Connect(config.Server)
	if config.Debug {
		conn.Debug = true
	}

	conn.Loop()

}
