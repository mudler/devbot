package team

import (
	"strings"

	"github.com/mudler/anagent"
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

func (m *Plugin) Register(a *anagent.Anagent) {
	log.Println("[Team] Started")
	if atom, err := bot.DBListKeys("timer"); err == nil {
		for _, s := range atom {
			v, err := bot.DBGetSingleKey("timer", s)
			if err == nil {
				SetupTimer(v, a)
			}
		}
	}
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

	bot.HandleCommand(msg, "timer list", func(args string) {
		if atom, err := bot.DBListKeys("timer"); err == nil {
			for _, t := range atom {
				if v, err := bot.DBGetSingleKey("timer", t); err == nil {
					conn.Privmsg(destination, t+": "+v)
				}
			}
		}
	})
	bot.HandleCommand(msg, "timer show", func(args string) {
		if v, err := bot.DBGetSingleKey("timer", args); err == nil {
			conn.Privmsg(destination, args+": "+v)
		}
	})

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

	bot.HandleCommand(msg, "timer add", func(args string) {
		split := strings.Split(args, " ")
		if ok := bot.DBPutKeyValue("timer", split[0], args); ok {
			conn.Privmsg(destination, "Added timer "+split[0]+": "+strings.Join(split, " "))
		}
		SetupTimer(args, bot.Agent)
	})

	bot.HandleCommand(msg, "timer remove", func(args string) {
		split := strings.Split(args, " ")
		bot.Agent.RemoveTimer(anagent.TimerID(split[0]))
		if ok := bot.DBRemoveSingleKey("timer", split[0]); ok {
			conn.Privmsg(destination, "Removed timer "+split[0])
		}
	})

}
