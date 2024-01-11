package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Settings     Settings     `yaml:"Settings,omitempty"`
	Applications Applications `yaml:"Application,omitempty"`
	Theme        Theme        `yaml:"Theme,omitempty"`
}

type Settings struct {
	HideHelpOnStartup bool `yaml:"HideHelpOnStartup"`
}

type Theme struct {
	PrimaryColor         string `yaml:"PrimaryColor"`
	ShadedColor          string `yaml:"ShadedColor"`
	HighlightColor       string `yaml:"HighlightColor"`
	HighlightShadedColor string `yaml:"HighlightShadedColor"`
	WarningColor         string `yaml:"WarningColor"`
	WarningShadedColor   string `yaml:"WarningShadedColor"`
	ReaderHighlightColor string `yaml:"ReaderHighlightColor"`
	ReaderHeadingColor   string `yaml:"ReaderHeadingColor"`
}

type Applications struct {
	Image Application `yaml:"Image,omitempty"`
	Audio Application `yaml:"Audio,omitempty"`
	Video Application `yaml:"Video,omitempty"`
	HTML  Application `yaml:"HTML,omitempty"`
}

type Application struct {
	Path string   `yaml:"Path"`
	Args []string `yaml:"Args"`
}

func Load(configFile string) (Configuration, error) {
	config := defaultConfiguration()
	// no config file supplied, use default values
	if configFile == "" {
		return config, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return Configuration{}, fmt.Errorf("Configuration error: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Configuration{}, fmt.Errorf("Configuration error: %s", err)
	}

	return config, nil
}

func defaultConfiguration() Configuration {
	return Configuration{
		Settings: Settings{
			HideHelpOnStartup: false,
		},
		Applications: Applications{},
		Theme:        gruvboxTheme(),
	}
}

func gruvboxTheme() Theme {
	return Theme{
		PrimaryColor:         "#EBDBB2",
		ShadedColor:          "#928374",
		HighlightColor:       "#458588",
		HighlightShadedColor: "#83A598",
		WarningColor:         "#FB4934",
		WarningShadedColor:   "#CC241D",
		ReaderHighlightColor: "#FABD2F",
		ReaderHeadingColor:   "#8EC07C",
	}
}
