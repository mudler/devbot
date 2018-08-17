package fun

import (
	"strings"

	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	util "github.com/mudler/devbot/shared/utils"
	irc "github.com/thoj/go-ircevent"

	"log"
)

type Slap struct{}

func init() {
	bot.RegisterPlugin(&Slap{})
}

func (m *Slap) Register(a *anagent.Anagent) {
	log.Println("[Slap] Started")
}

func (m *Slap) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	msg := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	if strings.Contains(msg, bot.Config.CommandPrefix+"slap") {
		args := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "slap")
		conn.Action(destination, " Slaps "+args+" with a large trout")
	}

}
