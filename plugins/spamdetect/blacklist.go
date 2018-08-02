package spamdetect

import (
	"strings"

	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	util "github.com/mudler/devbot/shared/utils"
	irc "github.com/thoj/go-ircevent"

	"log"
)

const BlacklistBucket = "blacklist"

type Blacklist struct{}

func init() {
	bot.RegisterPlugin(&Blacklist{})
}

func (m *Blacklist) Register(a *anagent.Anagent) {
	log.Println("[Blacklist] Started")
}

func (m *Blacklist) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	msg := event.Message()

	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}
	if bot.Config.IsAdmin(event.Nick) == true {
		if strings.Contains(msg, bot.Config.CommandPrefix+"blacklists") {
			if atom, err := bot.DBListKeys(BlacklistBucket); err == nil {
				for _, a := range atom {
					conn.Privmsg(destination, " - "+a)
				}
			}
		}

		if strings.Contains(msg, bot.Config.CommandPrefix+"blacklist add") {
			args := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "blacklist add")
			if bot.DBPutKey(BlacklistBucket, args) == true {
				conn.Privmsg(destination, args+" added to bad words blacklist")
			}
		}

		if strings.Contains(msg, bot.Config.CommandPrefix+"blacklist remove") {
			args := util.StripPluginCommand(msg, bot.Config.CommandPrefix, "blacklist remove")
			if bot.DBRemoveSingleKey(BlacklistBucket, args) == true {
				conn.Privmsg(destination, args+" removed from the blacklist")
			} else {
				conn.Privmsg(destination, args+" not found in the blacklist")
			}

		}
		return
	}

	if atom, err := bot.DBListKeys(BlacklistBucket); err == nil {
		for _, a := range atom {
			if strings.Contains(msg, a) {
				conn.SendRaw("MODE " + event.Arguments[0] + " +b " + event.Nick + "!*@*")
				conn.SendRaw("KICK " + event.Arguments[0] + " " + event.Nick + " : such topics are not very liked here.. RESOLVED->SPAM")
				return
			}
		}
	}

}
