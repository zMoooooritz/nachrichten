package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	AppConfig   ApplicationsConfig `yaml:"Application,omitempty"`
	ThemeConfig ThemeConfig        `yaml:"Theme,omitempty"`
}

type ThemeConfig struct {
	PrimaryColor           string `yaml:"PrimaryColor"`
	SecondaryColor         string `yaml:"SecondaryColor"`
	NormalTitleColor       string `yaml:"NormaleTitleColor"`
	NormalDescColor        string `yaml:"NormalDescColor"`
	SelectedPrimaryColor   string `yaml:"SelectedPrimaryColor"`
	SelectedSecondaryColor string `yaml:"SelectedSecondaryColor"`
	BreakingColor          string `yaml:"BreakingColor"`
}

type ApplicationsConfig struct {
	Image ApplicationConfig `yaml:"Image,omitempty"`
	Audio ApplicationConfig `yaml:"Audio,omitempty"`
	Video ApplicationConfig `yaml:"Video,omitempty"`
	HTML  ApplicationConfig `yaml:"HTML,omitempty"`
}

type ApplicationConfig struct {
	Path string   `yaml:"Path"`
	Args []string `yaml:"Args"`
}

type ResourceType int

const (
	TypeImage ResourceType = iota
	TypeAudio
	TypeVideo
	TypeHTML
)

type Config struct {
	file string
}

func Init() Config {
	cfg := Config{}

	cfg.file = *flag.String("config", "~/.config/nachrichten/config.yaml", "Path to configuration file")

	flag.Parse()
	return cfg
}

func (c Config) Load() Configuration {
	var config Configuration

	data, err := os.ReadFile(c.file)
	if err != nil {
		fmt.Println(err)
		return config
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}
