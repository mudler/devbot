package admin

import (
	"github.com/mudler/devbot/shared/registry"
	"github.com/mudler/devbot/shared/utils"
	"github.com/thoj/go-ircevent"

	"log"
	"strings"
)

type AdminPlugin struct{}

func init() {
	plugin_registry.RegisterPlugin(&AdminPlugin{})
}

func (m *AdminPlugin) Register() {
	log.Println("[AdminPlugin] Started")
}

func (m *AdminPlugin) OnPrivmsg(event *irc.Event) {
	conn := plugin_registry.Conn
	config := plugin_registry.Config
	message := event.Message()
	if config.IsAdmin(event.Nick) == false {
		return
	}

	if message == config.CommandPrefix+"help" {
		conn.Privmsg(event.Arguments[0], "- Admin commands - ")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"enable <plugin> - Load a specific plugin again in memory")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"disable <plugin> - UnLoad a specific plugin in memory")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"listplugins - List all plugins")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"op <nick> - Op nick on channel")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"deop <nick> - Deop nick on channel")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"kick <nick> - kick nick on channel")
	}

	if message == config.CommandPrefix+"listplugins" {
		ListPlugins(event.Arguments[0], conn)
	}
	if strings.Contains(message, config.CommandPrefix+"enable") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "enable")
		if plugin_registry.EnablePlugin(args) {
			conn.Privmsg(event.Arguments[0], args+" Enabled")
		}
		ListPlugins(event.Arguments[0], conn)
	}
	if strings.Contains(message, config.CommandPrefix+"disable") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "disable")
		if plugin_registry.DisablePlugin(args) {
			conn.Privmsg(event.Arguments[0], args+" Disabled")

		}
		ListPlugins(event.Arguments[0], conn)
	}

	if strings.Contains(message, config.CommandPrefix+"op") {
		cmd := strings.Split(message, " ")
		if cmd[1] != "" {
			log.Println("[AdminPlugin] MODE +o " + event.Arguments[0] + " " + cmd[1])
			conn.SendRaw("MODE +o " + event.Arguments[0] + " " + cmd[1])
		}
	}

	if strings.Contains(message, config.CommandPrefix+"deop") {
		cmd := strings.Split(message, " ")
		if cmd[1] != "" {
			log.Println("[AdminPlugin] MODE -o " + event.Arguments[0] + " " + cmd[1])

			conn.SendRaw("MODE -o " + event.Arguments[0] + " " + cmd[1])
		}
	}

	if strings.Contains(message, config.CommandPrefix+"kick") {
		cmd := strings.Split(message, " ")
		if cmd[1] != "" {
			log.Println("[AdminPlugin] KICK " + event.Arguments[0] + " " + cmd[1] + " : RESOLVED->KICKED")
			conn.SendRaw("KICK " + event.Arguments[0] + " " + cmd[1] + " : RESOLVED->KICKED")
		}
	}

}

func ListPlugins(sendTo string, conn *irc.Connection) {

	conn.Privmsg(sendTo, "Enabled plugins: ")
	for k, _ := range plugin_registry.Plugins {
		conn.Privmsg(sendTo, "\t"+k)

	}
	conn.Privmsg(sendTo, "Disabled plugins: ")

	for k, _ := range plugin_registry.DisabledPlugins {
		conn.Privmsg(sendTo, "\t"+k)
	}

}
