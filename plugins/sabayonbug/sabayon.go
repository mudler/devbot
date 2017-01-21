package sabayonbug

import (
	"github.com/mudler/devbot/plugins/gentoobug"
	"github.com/mudler/devbot/shared/registry"

	"github.com/thoj/go-ircevent"

	"log"
	"regexp"
)

type SabayonBugzillaPlugin struct{}

func init() {
	plugin_registry.RegisterPlugin(&SabayonBugzillaPlugin{})
}

func (m *SabayonBugzillaPlugin) Register() {
	log.Println("[SabayonBugzillaPlugin] Started")
}

func (m *SabayonBugzillaPlugin) OnPrivmsg(event *irc.Event) {
	conn := plugin_registry.Conn
	destination := event.Arguments[0]
	if event.Arguments[0] == plugin_registry.Config.BotNick {
		destination = event.Nick
	}

	// Detect if in chats are written bugs id like #12345
	regex, _ := regexp.Compile(`(?:^|\s)[＃#]{1}(\w+)`)
	bug := regex.FindStringSubmatch(event.Message())

	if len(bug) > 1 {
		buginfo := gentoobug.BugInfo("https://bugs.sabayon.org/show_bug.cgi?id=", bug[1])
		if buginfo.Summary != "" {
			conn.Privmsg(destination, buginfo.Url+"; "+buginfo.Summary+"; "+buginfo.AssignedTo+"; "+buginfo.Status+"; ")
		}
	}

}
