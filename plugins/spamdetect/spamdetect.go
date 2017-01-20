package spamdetect

import (
	"github.com/mudler/devbot/shared/registry"
	"github.com/mudler/devbot/shared/utils"

	"github.com/thoj/go-ircevent"
	"log"
)

type SpamDetectPlugin struct{}

func init() {
	plugin_registry.RegisterPlugin(&SpamDetectPlugin{})
}

func (m *SpamDetectPlugin) Register() {
	log.Println("[SpamDetectPlugin] Started")
}

func (m *SpamDetectPlugin) OnPrivmsg(event *irc.Event) {
	conn := plugin_registry.Conn
	config := plugin_registry.Config
	client := NewClient(config.UClassifyKey)
	results, _, _ := client.ClassifyService.Classify(config.UClassifyUser, config.UClassifyClassifier, []string{event.Message()})

	for i := range results {

		for c := range results[i].Classification {
			if results[i].Classification[c].ClassName == "unlegitimate" && results[i].Classification[c].P > 0.5 {
				conn.SendRaw("KICK " + event.Arguments[0] + " " + event.Nick + " : such topics are not very liked here.. RESOLVED->SPAM")
			}
			log.Println("[SpamDetectPlugin] Message: " + event.Message())
			log.Println("[SpamDetectPlugin] Classification: " + results[i].Classification[c].ClassName + " P: " + util.FloatToString(results[i].Classification[c].P))
			//conn.Privmsg(event.Arguments[0], " Classification: "+results[i].Classification[c].ClassName+" P: "+util.FloatToString(results[i].Classification[c].P))

		}

	}

}
