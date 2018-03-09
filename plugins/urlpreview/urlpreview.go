package urlpreview

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	"github.com/thoj/go-ircevent"
)

type UrlPreviewPlugin struct{}

func init() {
	bot.RegisterPlugin(&UrlPreviewPlugin{})
}

func (m *UrlPreviewPlugin) Register(a *anagent.Anagent) {
	log.Println("[UrlPreviewPlugin] Started")
}

func (m *UrlPreviewPlugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	message := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == bot.Config.BotNick {
		destination = event.Nick
	}

	if strings.Contains(message, "http://") || strings.Contains(message, "https://") || strings.Contains(message, "www.") {
		conn.Privmsg(destination, UrlTitle(message))
	}

}

// UrlTitle attempts to extract the title of the page that a
// pasted URL points to.
// Returns a string message with the title and URL on success, or a string
// with an error message on failure.
func UrlTitle(msg string) string {
	var (
		newMsg, url, title, word string
	)

	regex, _ := regexp.Compile(`(?i)<title>(.*?)<\/title>`)

	msgArray := strings.Split(msg, " ")

	for _, word = range msgArray {
		if strings.Contains(word, "http") {
			url = word
			break
		}
		if !strings.Contains(word, "http") && strings.Contains(word, "www") {
			url = "http://" + word
			break
		}
	}

	resp, err := http.Get(url)

	if err != nil {
		return fmt.Sprintf("Could not resolve URL %v, beware...\n", url)
	}

	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Could not read response Body of %v ...\n", url)
	}

	body := string(rawBody)
	noNewLines := strings.Replace(body, "\n", "", -1)
	noCarriageReturns := strings.Replace(noNewLines, "\r", "", -1)
	notSoRawBody := noCarriageReturns

	titleMatch := regex.FindStringSubmatch(notSoRawBody)
	if len(titleMatch) > 1 {
		title = strings.TrimSpace(titleMatch[1])
	} else {
		title = fmt.Sprintf("Title Resolution Failure")
	}
	newMsg = fmt.Sprintf("[ %v ]( %v )\n", title, url)

	return newMsg
}
