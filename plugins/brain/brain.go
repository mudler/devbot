package brain

import (
	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
	cobe "github.com/mudler/go.cobe"

	"log"
	"regexp"
	"strings"

	"github.com/thoj/go-ircevent"
)

type BrainPlugin struct {
	Brain *cobe.Cobe2Brain
}

func init() {
	bot.RegisterPlugin(&BrainPlugin{})
}

func (m *BrainPlugin) Register(a *anagent.Anagent) {
	log.Println("[BrainPlugin] Started")
	b, err := cobe.OpenCobe2Brain(bot.Config.BrainFile)
	m.Brain = b
	if err != nil {
		log.Fatalf("Opening brain file: %s", err)
	}
}

func (m *BrainPlugin) OnPrivmsg(event *irc.Event) {
	conn := bot.Conn
	config := bot.Config
	message := event.Message()
	destination := event.Arguments[0]
	if event.Arguments[0] == config.BotNick {
		destination = event.Nick
	}
	// message starts with command prefix, ignoring.
	if string(message[0]) == config.CommandPrefix {
		log.Println("[BrainPlugin] Do not learn !: " + message)
		return
	}

	re, _ := regexp.Compile(config.BotNick + "[:]?")
	sanitizedInput := re.ReplaceAllLiteralString(message, "")
	log.Println("[BrainPlugin] Learn: " + sanitizedInput)
	m.Brain.Learn(sanitizedInput)
	if !strings.HasPrefix(message, config.CommandPrefix) {
		if sanitizedInput != message || event.Arguments[0] == config.BotNick {
			log.Println("[BrainPlugin] Answer!")
			conn.Privmsg(destination, m.Brain.Reply(sanitizedInput))
		}
	}
}
