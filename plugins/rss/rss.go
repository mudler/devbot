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

	cmdArray := strings.SplitAfterN(msg, bot.Config.CommandPrefix, 2)
	if strings.Contains(msg, "subscribe") {

	}

}
