package seen

import (
	"fmt"
	"strings"
	"time"

	"github.com/mudler/devbot/bot"
	util "github.com/mudler/devbot/shared/utils"

	"github.com/thoj/go-ircevent"

	"log"
)

const timeFormat string = "2006-01-02 15:04:05 -0700 MST"

type Plugin struct {
}

func init() {
	bot.RegisterPlugin(&Plugin{})
}

func (m *Plugin) Register() {
	log.Println("[Team] Started")
}

func (m *Plugin) OnJoin(event *irc.Event) {
	bot.DBPutKeyValue("seen", event.Nick, time.Now().Format(timeFormat))
}

func (m *Plugin) OnQuit(event *irc.Event) {
	bot.DBPutKeyValue("seen", event.Nick, time.Now().Format(timeFormat))
}

func (m *Plugin) OnPart(event *irc.Event) {
	bot.DBPutKeyValue("seen", event.Nick, time.Now().Format(timeFormat))
}

func (m *Plugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	msg := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	bot.DBPutKeyValue("seen", event.Nick, time.Now().Format(timeFormat))

	if strings.Contains(msg, bot.Config.CommandPrefix+"seen") {
		args := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "seen")
		if seen, err := bot.DBGetSingleKey("seen", args); err == nil {
			if len(seen) > 0 {
				then, err := time.Parse(timeFormat, seen)
				if err != nil {
					conn.Privmsg(destination, "Error parsing time: "+err.Error())
					return
				}
				duration := time.Since(then)
				conn.Privmsg(destination, args+" not seen since "+fmt.Sprint(duration.Minutes())+"m"+" ("+seen+")")
			} else {
				conn.Privmsg(destination, "User never seen")
			}
		} else {
			conn.Privmsg(destination, "Error while retrieving the data")
		}

	}

}
