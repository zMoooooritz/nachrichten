package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	SettingsConfig SettingsConfig     `yaml:"Settings"`
	AppConfig      ApplicationsConfig `yaml:"Application,omitempty"`
	Theme          Theme              `yaml:"Theme,omitempty"`
}

type SettingsConfig struct {
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

func Load(configFile string) Configuration {
	config := defaultConfiguration()

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config
	}

	_ = yaml.Unmarshal(data, &config)
	return config
}

func defaultConfiguration() Configuration {
	return Configuration{
		SettingsConfig: SettingsConfig{
			HideHelpOnStartup: false,
		},
	}
}
