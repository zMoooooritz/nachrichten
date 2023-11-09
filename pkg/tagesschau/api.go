package tagesschau

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"github.com/zMoooooritz/Nachrichten/pkg/http"
)

const (
	baseUrl      string = "https://www.tagesschau.de/"
	apiUrl       string = baseUrl + "api2u/homepage"
	shortNewsUrl string = baseUrl + "multimedia/sendung/tagesschau_in_100_sekunden"
)

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

func (n NewsEntry) Title() string       { return n.Topline }
func (n NewsEntry) Description() string { return n.Desc }
func (n NewsEntry) FilterValue() string { return n.Topline }

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

func LoadNews() News {
	body, err := http.FetchURL(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	var news News
	err = json.Unmarshal(body, &news)
	if err != nil {
		log.Fatal(err)
	}
	return news
}

func GetShortNewsURL() (string, error) {
	body, err := http.FetchURL(shortNewsUrl)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	data, exists := doc.Find("div.teaser__media").Find("div.v-instance").Attr("data-v")
	if !exists {
		return "", errors.New("Unable to parse HTML to find URL")
	}

	url, err := jsonparser.GetString([]byte(data), "mc", "streams", "[0]", "media", "[4]", "url")
	if err != nil {
		return "", err
	}
	return url, nil
}
