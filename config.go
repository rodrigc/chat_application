package chat_application

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ChatConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"telnet_port"`
	RestAPIPort int    `yaml:"rest_api_port"`
	LogFile     string `yaml:"logfile"`
}

var chatConfig ChatConfig

// ParseConfigFile will read a yaml file which
// contains configuration information for the chat server.
// See chat.yml for the format of the yaml file.
func ParseConfigFile() error {
	// TODO: for now hardcode the config file as
	// chat.yml.  This config file must be in the same
	// directory as the chatserver binary.
	// Make this configurable later.
	configFileName := "chat.yml"

	configFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFile, &chatConfig)
	if err != nil {
		return err
	}

	return nil
}
