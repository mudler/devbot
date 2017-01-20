package util

import (
	"encoding/json"
	"fmt"
	"os"
)

//Configuration needed for plugins and bot
type Config struct {
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
	Debug               bool
	Welcome             bool
	WelcomeMessage      string
	JoinMessage         string
	MessageOnJoin       bool
	Tls                 bool
	Administrators      map[string]bool
	Plugins             map[string]bool
}

func (c *Config) IsAdmin(user string) bool {
	_, ok := c.Administrators[user]
	return ok
}

func (c *Config) EnabledPlugin(plugin string) bool {
	v, _ := c.Plugins[plugin]
	return v
}

func LoadConfig(f string) (Config, error) {

	file, err := os.Open(f)

	if err != nil {
		fmt.Println("Couldn't read config file")
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	config.BrainFile = "cobe.brain"
	config.MessageOnJoin = false
	config.Welcome = false
	config.Debug = false
	config.UClassifyClassifier = "Spam"
	config.UClassifyUser = "mudler"
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Couldn't parse json file")
		return Config{}, err
	}
	return config, err

}
