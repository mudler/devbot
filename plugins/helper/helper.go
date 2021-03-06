package helper

import (
	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	"github.com/thoj/go-ircevent"

	"log"
	"strings"
)

type HelperPlugin struct{}

func init() {
	bot.RegisterPlugin(&HelperPlugin{})
}

func (m *HelperPlugin) Register(a *anagent.Anagent) {
	log.Println("[HelperPlugin] Started")
}

func (m *HelperPlugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	config := bot.Config
	message := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == config.BotNick {
		destination = event.Nick
	}

	switch {
	case strings.Contains(message, "help"):
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"wiki - Display wiki url")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"homepage - Display homepage url")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"forum - Display forum url")
		conn.Privmsg(destination, "\t"+config.CommandPrefix+"bugs - Display bugzilla url")
	case strings.Contains(message, "wiki"):
		conn.Privmsg(destination, config.WikiLink)
	case strings.Contains(message, "homepage"):
		conn.Privmsg(destination, config.Homepage)
	case strings.Contains(message, "forum"):
		conn.Privmsg(destination, config.Forums)
	case strings.Contains(message, "bugs"):
		conn.Privmsg(destination, config.Bugs)

	}

}
