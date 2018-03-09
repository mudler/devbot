package team

import (
	"strings"

	"github.com/mudler/devbot/bot"
	util "github.com/mudler/devbot/shared/utils"

	"github.com/thoj/go-ircevent"

	"log"
)

type Plugin struct {
}

func init() {
	bot.RegisterPlugin(&Plugin{})
}

func (m *Plugin) Register() {
	log.Println("[Team] Started")
}

func (m *Plugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	msg := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	if strings.Contains(msg, bot.Config.CommandPrefix+"team") {
		args := strings.TrimSpace(util.StripPluginCommand(msg, bot.Config.CommandPrefix, "team"))
		if atom, err := bot.DBListKeys("team" + args); err == nil {
			conn.Privmsg(destination, args+": "+strings.Join(atom, " "))
		}
	}

	if bot.Config.IsAdmin(event.Nick) == false {
		return
	}

	if strings.Contains(msg, bot.Config.CommandPrefix+"member add") {
		args := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "member add")
		split := strings.Split(args, " ")
		if len(split) > 1 {
			if bot.DBPutKey("team"+split[0], split[1]) == true {
				conn.Privmsg(destination, split[1]+" now is member of "+split[0])
			}
		}
	}

	if strings.Contains(msg, bot.Config.CommandPrefix+"member remove") {
		args := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "member remove")
		split := strings.Split(args, " ")
		if len(split) > 1 {
			if bot.DBRemoveSingleKey("team"+split[0], split[1]) == true {
				conn.Privmsg(destination, split[1]+" now is not anymore member of "+split[0])
			}
		}
	}

}
