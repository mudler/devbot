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
  destination := event.Arguments[0]
  if event.Arguments[0] == config.BotNick {
    destination = event.Nick
  }
  // message starts with command prefix, ignoring.
  if string(message[0]) == config.CommandPrefix {
    log.Println("[BrainPlugin] Do not learn !: "+message)
    return
  }

	re, _ := regexp.Compile(config.BotNick + "[:]?")
	sanitizedInput := re.ReplaceAllLiteralString(message, "")
  log.Println("[BrainPlugin] Learn: "+sanitizedInput)
	m.Brain.Learn(sanitizedInput)
	if !strings.HasPrefix(message, config.CommandPrefix) {
		if sanitizedInput != message || event.Arguments[0] == config.BotNick {
      log.Println("[BrainPlugin] Answer!")
			conn.Privmsg(destination, m.Brain.Reply(sanitizedInput))
		}
	}
}
