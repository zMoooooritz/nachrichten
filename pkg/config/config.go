package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	SettingsConfig SettingsConfig     `yaml:"Settings"`
	AppConfig      ApplicationsConfig `yaml:"Application,omitempty"`
	ThemeConfig    ThemeConfig        `yaml:"Theme,omitempty"`
}

type SettingsConfig struct {
	HideHelpOnStartup bool `yaml:"HideHelpOnStartup"`
}

type ThemeConfig struct {
	PrimaryColor           string `yaml:"PrimaryColor"`
	SecondaryColor         string `yaml:"SecondaryColor"`
	NormalTitleColor       string `yaml:"NormaleTitleColor"`
	NormalDescColor        string `yaml:"NormalDescColor"`
	SelectedPrimaryColor   string `yaml:"SelectedPrimaryColor"`
	SelectedSecondaryColor string `yaml:"SelectedSecondaryColor"`
	BreakingColor          string `yaml:"BreakingColor"`
	ReaderHighlightColor   string `yaml:"ReaderHighlightColor"`
	ReaderHeadingColor     string `yaml:"ReaderHeadingColor"`
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
	ConfigFile string
	LogFile    string
}

func Init() Config {

	configFile := flag.String("config", "", "Path to configuration file")
	logFile := flag.String("debug", "", "Path to log file")

	flag.Parse()
	cfg := Config{
		ConfigFile: *configFile,
		LogFile:    *logFile,
	}
	return cfg
}

func (c Config) Load() Configuration {
	config := defaultConfiguration()

	data, err := os.ReadFile(c.ConfigFile)
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
