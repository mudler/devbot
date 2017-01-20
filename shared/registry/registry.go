// Plugin register themselves here, the registry keep tracks of plugins to redirect the messages
package plugin_registry

import (
	"github.com/mudler/devbot/shared/utils"
	"github.com/thoj/go-ircevent"
	"log"
	"reflect"
	"strings"
)

type DevbotPlugin interface {
	Register()
}

// These are are registered plugins
var Plugins = map[string]DevbotPlugin{}
var DisabledPlugins = map[string]DevbotPlugin{}
var Commands = make(map[string]string)
var Config util.Config
var Conn *irc.Connection

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
