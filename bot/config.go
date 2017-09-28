package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

//Configurationuration needed for plugins and bot
type Configuration struct {
	// !quit !part the prefix is "!"
	CommandPrefix string
	// Server to connect to, format ip:port
	Server string
	// Separate channel list to join
	Channel             []string
	BotUser             string
	BotNick             string
	BotNickPassword     string
	Trigger             string
	WeatherKey          string
	UClassifyKey        string
	UClassifyUser       string
	UClassifyClassifier string
	LogDir              string
	WikiLink            string
	Homepage            string
	Forums              string
	Bugs                string
	BrainFile           string
	DBFile              string
	Debug               bool
	Welcome             bool
	WelcomeMessage      string
	JoinMessage         string
	MessageOnJoin       bool
	Tls                 bool
	Administrators      map[string]bool
	Plugins             map[string]bool
}

func (c *Configuration) IsAdmin(user string) bool {
	_, ok := c.Administrators[user]
	return ok
}

func (c *Configuration) EnabledPlugin(plugin string) bool {
	v, _ := c.Plugins[plugin]
	return v
}

func LoadConfig(f string) (Configuration, error) {

	file, err := os.Open(f)

	if err != nil {
		fmt.Println("Couldn't read Configuration file")
		return Configuration{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var Config Configuration
	Config.BrainFile = "cobe.brain"
	Config.DBFile = "bot.db"
	Config.MessageOnJoin = false
	Config.Welcome = false
	Config.Debug = false
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Couldn't parse json file")
		return Configuration{}, err
	}
	return Config, err
}
