// Plugin register themselves here, the registry keep tracks of plugins to redirect the messages
package bot

import (
	"crypto/tls"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/thoj/go-ircevent"
)

type DevbotPlugin interface {
	Register()
}

// These are are registered plugins
var Plugins = map[string]DevbotPlugin{}
var DisabledPlugins = map[string]DevbotPlugin{}
var Commands = make(map[string]string)
var Config Configuration
var Conn *irc.Connection
var DB string

// Register a Plugin
func RegisterPlugin(p DevbotPlugin) {
	Plugins[KeyOf(p)] = p
}

// Disable a plugin
func DisablePlugin(plugin string) bool {
	plugin = strings.TrimSpace(plugin)
	_, exists := Plugins[plugin]
	if exists {
		DisabledPlugins[plugin] = Plugins[plugin]
		_, disabled := DisabledPlugins[plugin]
		if disabled {
			delete(Plugins, plugin)
			//		DisabledPlugins[plugin].OnStop()

			log.Println(plugin + " removed from running plugins")
		} else {
			log.Println("Can't disable " + plugin + ", odd")

		}
		return disabled
	} else {
		log.Println("Plugin '" + plugin + "' does not exist or is not loaded")

	}
	return exists
}

// Enable a plugin
func EnablePlugin(plugin string) bool {
	plugin = strings.TrimSpace(plugin)

	_, PluginExists := Plugins[plugin]
	if PluginExists {
		return true
	}

	PluginInstance, InstanceExists := DisabledPlugins[plugin]
	Plugins[plugin] = PluginInstance
	if InstanceExists {

		delete(DisabledPlugins, plugin)
		//PluginInstance.OnStart()

		log.Println(plugin + " enabled ")
		return true
	}
	return false
}

func KeyOf(p DevbotPlugin) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}

// Register a Command exported by a plugin
func RegisterCommand(command string, description string) {
	Commands[command] = description
}

// UnRegister a Command exported by a plugin
func UnregisterCommand(command string) {
	delete(Commands, command)
}

func Start(config Configuration) {
	conn := irc.IRC(config.BotNick, config.BotUser)

	if config.UnsecureTLS == true {
		conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	conn.UseTLS = config.Tls
	fmt.Println("BotNick:\t" + config.BotNick)
	fmt.Println("BotUser:\t" + config.BotUser)
	fmt.Println("Channels:")
	for _, channel := range config.Channel {
		fmt.Println("\t" + channel)
	}

	Config = config
	Conn = conn
	DB = config.DBFile

	// Bootstrapper for plugins
	for _, d := range Plugins {
		if (len(config.Plugins) > 0 && config.EnabledPlugin(KeyOf(d))) || len(config.Plugins) == 0 {
			go d.Register()
		}
	}

	// Bootstrapper for plugins
	for _, d := range Plugins {
		if len(config.Plugins) > 0 && !config.EnabledPlugin(KeyOf(d)) {
			DisablePlugin(KeyOf(d))
		}
	}

	log.Println(strconv.Itoa(len(Plugins)) + " plugins loaded")
	RegisterCallbacks(conn)
	if err := conn.Connect(config.Server); err != nil {
		log.Println("Connection failed: " + err.Error())
		return
	}
	if config.Debug {
		conn.Debug = true
	}
	conn.Loop()
}

func RegisterCallbacks(conn *irc.Connection) {

	conn.AddCallback("001", func(e *irc.Event) {
		for _, channel := range Config.Channel {
			conn.Join(channel)
		}
		if Config.BotNickPassword != "" {
			log.Println("Identifying bot nickname against NickServ")
			conn.Privmsg("NickServ", "IDENTIFY "+Config.BotNickPassword)
		}
	})

	conn.AddCallback("PRIVMSG", func(event *irc.Event) {
		for _, d := range Plugins {
			if obj, ok := d.(interface {
				OnPrivmsg(*irc.Event)
			}); ok {
				go obj.OnPrivmsg(event)
			}
		}
	})
	conn.AddCallback("JOIN", func(event *irc.Event) {
		for _, d := range Plugins {
			if obj, ok := d.(interface {
				OnJoin(*irc.Event)
			}); ok {
				go obj.OnJoin(event)
			}
		}

	})
	conn.AddCallback("PART", func(event *irc.Event) {
		for _, d := range Plugins {
			if obj, ok := d.(interface {
				OnPart(*irc.Event)
			}); ok { //Check if plugin has that method
				go obj.OnPart(event)
			}
		}
	})
	conn.AddCallback("QUIT", func(event *irc.Event) {
		for _, d := range Plugins {
			if obj, ok := d.(interface {
				OnQuit(*irc.Event)
			}); ok { //Check if plugin has that method
				go obj.OnQuit(event)
			}
		}

	})
}
