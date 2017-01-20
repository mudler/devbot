package gentoobug

import (
	"github.com/mudler/devbot/shared/registry"
	"github.com/thoj/go-ircevent"
	"log"
	"regexp"
)

type GentooBugzillaPlugin struct{}

func init() {
	plugin_registry.RegisterPlugin(&GentooBugzillaPlugin{})
}

func (m *GentooBugzillaPlugin) Register() {
	log.Println("[GentooBugzillaPlugin] Started")
}

func (m *GentooBugzillaPlugin) OnPrivmsg(event *irc.Event) {

	conn := plugin_registry.Conn

	// Detect if in chats are written bugs id like #12345
	regex, _ := regexp.Compile(`(?:^|\s)[＃#]{1}(\w+)`)

	bug := regex.FindStringSubmatch(event.Message())

	if len(bug) > 1 {
		buginfo := BugInfo("https://bugs.gentoo.org/show_bug.cgi?id=", bug[1])
		if buginfo.Summary != "" {
			conn.Privmsg(event.Arguments[0], buginfo.Url+"; "+buginfo.Summary+"; "+buginfo.AssignedTo+"; "+buginfo.Status+"; ")
		}
	}

}
