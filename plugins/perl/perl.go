package perl

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	"github.com/thoj/go-ircevent"
)

type Plugin struct{}

func init() {
	bot.RegisterPlugin(&Plugin{})
}

func (m *Plugin) Register(a *anagent.Anagent) {
	log.Println("[Perl] Started")
}

func (m *Plugin) OnPrivmsg(event *irc.Event) {

	conn := bot.Conn
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	regex, _ := regexp.Compile(`(?i)` + bot.Config.CommandPrefix + `perl(?:|\s)[:]{1}(.*)`)

	code := regex.FindStringSubmatch(event.Message())

	if len(code) > 1 {
		file, err := ioutil.TempFile(os.TempDir(), "prefix")
		defer os.Remove(file.Name())
		file.WriteString(code[0])
		defer file.Close()
		file.Sync()
		out, err := exec.Command("/usr/bin/perl", "-T", file.Name()).CombinedOutput()
		if err != nil {
			conn.Privmsg(destination, "Error: "+err.Error())
		}
		conn.Privmsg(destination, "Output: "+string(out))
	}

}
