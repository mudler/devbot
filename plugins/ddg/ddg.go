package ddg

import (
	"fmt"

	"github.com/darthlukan/goduckgo/goduckgo"
	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	"github.com/mudler/devbot/plugins/urlpreview"

	"log"
	"strings"

	"github.com/thoj/go-ircevent"
)

type DDGPlugin struct{}

func init() {
	bot.RegisterPlugin(&DDGPlugin{})
}

func (m *DDGPlugin) Register(a *anagent.Anagent) {
	log.Println("[DDGPlugin] Started")
}

func (m *DDGPlugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	config := bot.Config
	var (
		msg      string
		msgArray []string
		cmdArray []string
	)
	msg = event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == config.BotNick {
		destination = event.Nick
	}

	cmdArray = strings.SplitAfterN(msg, config.CommandPrefix, 2)
	if !strings.Contains(msg, "ddg") && !strings.Contains(msg, "search") {
		return
	}
	if len(cmdArray) > 1 {
		msgArray = strings.SplitN(cmdArray[1], " ", 2)
	}

	if len(msgArray) > 1 {
		query := strings.Join(msgArray[1:], " ")
		conn.Privmsg(destination, SearchCmd(query))
	}

}

// WebSearch takes a query string as an argument and returns
// a formatted string containing the results from DuckDuckGo.
func SearchCmd(query string) string {
	msg, err := goduckgo.Query(query)
	if err != nil {
		return fmt.Sprintf("DDG Error: %v\n", err)
	}

	switch {
	case len(msg.RelatedTopics) > 0:
		return fmt.Sprintf("First Topical Result: [ %s ]( %s )\n", msg.RelatedTopics[0].FirstURL, msg.RelatedTopics[0].Text)
	case len(msg.Results) > 0:
		return fmt.Sprintf("First External result: [ %s ]( %s )\n", msg.Results[0].FirstURL, msg.Results[0].Text)
	case len(msg.Redirect) > 0:
		return fmt.Sprintf("Redirect result: %s\n", urlpreview.UrlTitle(msg.Redirect))
	default:
		return fmt.Sprintf("Query: '%s' returned no results.\n", query)
	}
}
