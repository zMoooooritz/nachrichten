package main

import "time"

type News struct {
	NationalNews []NewsEntry `json:"news"`
	RegionalNews []NewsEntry `json:"regional"`
}

type NewsEntry struct {
	Topline      string    `json:"topline"`
	Desc         string    `json:"title"`
	Introduction string    `json:"firstSentence"`
	URL          string    `json:"shareURL"`
	Breaking     bool      `json:"breakingNews"`
	Date         time.Time `json:"date"`
	Image        Image     `json:"teaserImage"`
	Content      []Content `json:"content"`
	Video        Video     `json:"video"`
}

type Content struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Image struct {
	Title     string    `json:"alttext"`
	Type      string    `json:"type"`
	ImageURLs ImageURLs `json:"imageVariants"`
}

type ImageURLs struct {
	SquareSmall  string `json:"1x1-144"`
	SquareMedium string `json:"1x1-432"`
	SquareBig    string `json:"1x1-840"`
	RectSmall    string `json:"16x9-256"`
	RectMedium   string `json:"16x9-640"`
	RectBig      string `json:"16x9-1920"`
}

type Video struct {
	Title     string    `json:"title"`
	VideoURLs VideoURLs `json:"streams"`
}

type VideoURLs struct {
	Small  string `json:"h264s"`
	Medium string `json:"h264m"`
	Big    string `json:"h264xl"`
}

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
