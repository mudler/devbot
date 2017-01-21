package logger

import (
	"github.com/mudler/devbot/shared/registry"
	"github.com/thoj/go-ircevent"

	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type LoggerPlugin struct{}

func init() {
	plugin_registry.RegisterPlugin(&LoggerPlugin{})
}

func (m *LoggerPlugin) Register() {
	log.Println("[LoggerPlugin] Started")
}
func (m *LoggerPlugin) OnQuit(event *irc.Event) {
	message := fmt.Sprintf("%s exited", event.Nick)
	go ChannelLogger(plugin_registry.Config.LogDir+"/"+event.Arguments[0], event.Nick, message)
}
func (m *LoggerPlugin) OnPart(event *irc.Event) {
	message := fmt.Sprintf("%s left the room", event.Nick)
	go ChannelLogger(plugin_registry.Config.LogDir+"/"+event.Arguments[0], event.Nick, message)
}
func (m *LoggerPlugin) OnPrivmsg(event *irc.Event) {
	destination := event.Arguments[0]
	if event.Arguments[0] == plugin_registry.Config.BotNick {
		destination = event.Nick
	}
	go ChannelLogger(plugin_registry.Config.LogDir+"/"+destination, event.Nick+": ", event.Message())
}

func (m *LoggerPlugin) OnJoin(event *irc.Event) {
	config := plugin_registry.Config

	if event.Nick == config.BotNick {
		LogDir(config.LogDir)
		LogFile(config.LogDir + "/" + event.Arguments[0])
	}
	message := fmt.Sprintf("%s has joined", event.Nick)
	go ChannelLogger(config.LogDir+"/"+event.Arguments[0], event.Nick, message)
}

func ChannelLogger(Log string, UserNick string, message string) {
	STime := time.Now().UTC().Format(time.ANSIC)
	log_file := strings.Replace(Log, "#", "", 1)
	logFile := fmt.Sprintf("%s.log", log_file)

	//Open the file for writing With Append Flag to create file persistence
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println(err)
		LogFile(log_file)
		f, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_SYNC, 0666)
	}

	defer f.Close()
	n, err := io.WriteString(f, fmt.Sprintf("%v > %v: %v\n", STime, UserNick, message))
	if err != nil {
		fmt.Println(n, err)
	}
}

func LogDir(CreateDir string) {
	if _, err := os.Stat(CreateDir); os.IsNotExist(err) {
		os.Mkdir(CreateDir, 0777)
	}
}

func LogFile(CreateFile string) {
	log := strings.Replace(CreateFile, "#", "", 1)
	logFile := fmt.Sprintf("%s.log", log)

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		os.Create(logFile)
	}
}
