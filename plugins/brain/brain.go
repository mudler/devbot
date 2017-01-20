package brain

import (
	"github.com/mudler/devbot/shared/registry"
	cobe "github.com/mudler/go.cobe"

	"github.com/thoj/go-ircevent"
	"log"
	"regexp"
	"strings"
)

type BrainPlugin struct {
	Brain *cobe.Cobe2Brain
}

func init() {
	plugin_registry.RegisterPlugin(&BrainPlugin{})
}

func (m *BrainPlugin) Register() {
	log.Println("[BrainPlugin] Started")
	b, err := cobe.OpenCobe2Brain(plugin_registry.Config.BrainFile)
	m.Brain = b
	if err != nil {
		log.Fatalf("Opening brain file: %s", err)
	}
}

func (m *BrainPlugin) OnPrivmsg(event *irc.Event) {
	conn := plugin_registry.Conn
	config := plugin_registry.Config
	message := event.Message()

	re, _ := regexp.Compile(config.BotNick + "[:]?")
	sanitizedInput := re.ReplaceAllLiteralString(message, "")
	m.Brain.Learn(sanitizedInput)
	if !strings.HasPrefix(message, config.CommandPrefix) {
		if sanitizedInput != message {
			conn.Privmsg(event.Arguments[0], m.Brain.Reply(sanitizedInput))
		}
	}
}
