package admin

import (
	"log"
	"strings"

	"github.com/mudler/anagent"
	irc "github.com/thoj/go-ircevent"

	"github.com/mudler/devbot/bot"
)

type AutoVoice struct{}

func init() {
	bot.RegisterPlugin(&AutoVoice{})
}

func (m *AutoVoice) Register(a *anagent.Anagent) {
	log.Println("[AutoVoice] Started")
}

func (m *AutoVoice) OnJoin(event *irc.Event) {
	conn := bot.Conn
	config := bot.Config

	//message := event.Message()
	channel := event.Arguments[0]

	if strings.HasPrefix(event.Nick, config.AutoVoicePrefix) {
		conn.SendRaw("MODE " + channel + " +v " + event.Nick)
	}

}
