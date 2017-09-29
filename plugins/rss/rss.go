package rss

import (
	"strings"
	"time"

	"github.com/mudler/devbot/bot"
	util "github.com/mudler/devbot/shared/utils"

	"github.com/thoj/go-ircevent"

	"log"
)

type Plugin struct {
	CheckFeedtimer chan bool
}

func init() {
	bot.RegisterPlugin(&Plugin{})
}

func (m *Plugin) Register() {
	log.Println("[RSSRead] Started")
	m.CheckFeedtimer = util.RecurringTimer(CheckFeed, 10*time.Second)
}

func (m *Plugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	msg := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	if bot.Config.IsAdmin(event.Nick) == false {
		return
	}

	if strings.Contains(msg, bot.Config.CommandPrefix+"subscribe") {
		atom_url := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "subscribe")
		bot.DBPutKey("subscriptions", atom_url)
		if bot.DBPutKey("atom"+atom_url, destination) == true {
			conn.Privmsg(destination, destination+" now is subscribed to "+atom_url)
		}
	}

	if strings.Contains(msg, bot.Config.CommandPrefix+"unsubscribe") {
		atom_url := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "unsubscribe")
		if bot.DBRemoveSingleKey("atom"+atom_url, destination) == true {
			conn.Privmsg(destination, destination+" now is unsubscribed from "+atom_url)
		}
	}

}
