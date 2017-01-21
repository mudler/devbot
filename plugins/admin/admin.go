package admin

import (
	"github.com/mudler/devbot/shared/registry"
	"github.com/mudler/devbot/shared/utils"
	"github.com/thoj/go-ircevent"
	"net/http"

	"github.com/inconshreveable/go-update"

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
	destination := event.Arguments[0]
	if event.Arguments[0] == config.BotNick {
		destination = event.Nick
	}
	if config.IsAdmin(event.Nick) == false {
		return
	}

	if message == config.CommandPrefix+"help" {
		conn.Privmsg(destination, "- Admin commands - ")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"enable <plugin> - Load a specific plugin again in memory")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"disable <plugin> - UnLoad a specific plugin in memory")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"listplugins - List all plugins")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"op <nick> - Op nick on channel")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"deop <nick> - Deop nick on channel")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"kick <nick> - kick nick on channel")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"join <channel> - join channel")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"update <url> - update bot with the given url")

	}

	if message == config.CommandPrefix+"listplugins" {
		ListPlugins(destination, conn)
	}
	if strings.Contains(message, config.CommandPrefix+"enable") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "enable")
		if plugin_registry.EnablePlugin(args) {
			conn.Privmsg(destination, args+" Enabled")
		}
		ListPlugins(destination, conn)
	}
	if strings.Contains(message, config.CommandPrefix+"disable") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "disable")
		if plugin_registry.DisablePlugin(args) {
			conn.Privmsg(destination, args+" Disabled")
		}
		ListPlugins(destination, conn)
	}

	if strings.Contains(message, config.CommandPrefix+"op") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "join")
		if args != "" {
			conn.SendRaw("MODE +o " + destination + " " + args)
		}
	}
	if strings.Contains(message, config.CommandPrefix+"join") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "join")
		if args != "" {
			conn.Join(args)
		}
	}
	if strings.Contains(message, config.CommandPrefix+"part") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "part")
		if args != "" {
			conn.SendRaw("PART " + args)
		}
	}

	if strings.Contains(message, config.CommandPrefix+"deop") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "part")
		if args != "" {
			conn.SendRaw("MODE -o " + destination + " " + args)
		}
	}

	if strings.Contains(message, config.CommandPrefix+"kick") {
		args := util.StripPluginCommand(message, config.CommandPrefix, "part")
		if args != "" {
			conn.SendRaw("KICK " + destination + " " + args + " :RESOLVED->KICKED")
		}
	}

	if strings.Contains(message, config.CommandPrefix+"update") {
		url := util.StripPluginCommand(message, config.CommandPrefix, "update")
		if url != "" {
			err := doUpdate(url)
			if err != nil {
				conn.Privmsg(destination, err.Error())
			}
		}
	}

}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	return err
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
