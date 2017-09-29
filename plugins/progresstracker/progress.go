package progresstracker

import (
	"log"
	"regexp"

	"github.com/mudler/devbot/bot"
	"github.com/thoj/go-ircevent"
)

type Plugin struct{}

func init() {
	bot.RegisterPlugin(&Plugin{})
}

func (m *Plugin) Register() {
	log.Println("[ProgressTracker] Started")
}

func (m *Plugin) OnPrivmsg(event *irc.Event) {

	conn := bot.Conn
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	// Detect if in chats are written bugs id like #12345
	regex, _ := regexp.Compile(`(?i)poo(?:|\s)[＃#]{1}(\w+)`)

	bug := regex.FindStringSubmatch(event.Message())

	if len(bug) > 1 {
		buginfo := BugInfo("https://progress.opensuse.org/issues/", bug[1])
		if buginfo.Summary != "" {
			conn.Privmsg(destination, buginfo.Url+"; "+buginfo.Summary+"; "+buginfo.Status+"; ")
		}
	}

}
