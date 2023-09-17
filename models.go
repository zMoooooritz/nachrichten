package main

import "time"

type News struct {
	NationalNews []NewsEntry `json:"news"`
	RegionalNews []NewsEntry `json:"regional"`
}

type NewsEntry struct {
	TopLine      string    `json:"topline"`
	Title        string    `json:"title"`
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
