package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BotToken          string `yaml:"bot_token"`
	CustomBotAPI      string `yaml:"custom_bot_api"`
	DropPendingUpdate bool   `yaml:"drop_pending_update"`
}

func ParseConfig(configFile string) (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(f)
	var config Config
	err = dec.Decode(&config)
	return &config, err
}
