package gentoobug

import (
	"log"
	"regexp"

	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	"github.com/thoj/go-ircevent"
)

type GentooBugzillaPlugin struct{}

func init() {
	bot.RegisterPlugin(&GentooBugzillaPlugin{})
}

func (m *GentooBugzillaPlugin) Register(a *anagent.Anagent) {
	log.Println("[GentooBugzillaPlugin] Started")
}

func (m *GentooBugzillaPlugin) OnPrivmsg(event *irc.Event) {

	conn := bot.Conn
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	// Detect if in chats are written bugs id like #12345
	regex, _ := regexp.Compile(`(?:^|\s)[＃#]{1}(\w+)`)

	bug := regex.FindStringSubmatch(event.Message())

	if len(bug) > 1 {
		buginfo := BugInfo("https://bugs.gentoo.org/show_bug.cgi?id=", bug[1])
		if buginfo.Summary != "" {
			conn.Privmsg(destination, buginfo.Url+" - "+buginfo.Summary+" - "+buginfo.AssignedTo+" - "+buginfo.Status)
		}
	}

}
