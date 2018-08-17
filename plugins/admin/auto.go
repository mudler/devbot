package admin

import (
	"log"
	"regexp"

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

	match, err := regexp.MatchString(config.AutoVoicePrefix, event.Nick)
	if err != nil {
		return
	}

	//message := event.Message()
	channel := event.Arguments[0]

	if match {
		conn.SendRaw("MODE " + channel + " +v " + event.Nick)
	}

}
