package helper

import (
	"github.com/mudler/devbot/shared/registry"
	"github.com/thoj/go-ircevent"

	"log"
	"strings"
)

type HelperPlugin struct{}

func init() {
	plugin_registry.RegisterPlugin(&HelperPlugin{})
}

func (m *HelperPlugin) Register() {
	log.Println("[HelperPlugin] Started")
}

func (m *HelperPlugin) OnPrivmsg(event *irc.Event) {
	conn := plugin_registry.Conn
	config := plugin_registry.Config
	message := event.Message()

	switch {
	case strings.Contains(message, "help"):
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"wiki - Display wiki url")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"homepage - Display homepage url")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"forum - Display forum url")
		conn.Privmsg(event.Arguments[0], "\t"+config.CommandPrefix+"bugs - Display bugzilla url")
	case strings.Contains(message, "wiki"):
		conn.Privmsg(event.Arguments[0], config.WikiLink)
	case strings.Contains(message, "homepage"):
		conn.Privmsg(event.Arguments[0], config.Homepage)
	case strings.Contains(message, "forum"):
		conn.Privmsg(event.Arguments[0], config.Forums)
	case strings.Contains(message, "bugs"):
		conn.Privmsg(event.Arguments[0], config.Bugs)

	}

}
