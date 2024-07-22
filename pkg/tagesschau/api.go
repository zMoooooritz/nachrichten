package tagesschau

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"github.com/zMoooooritz/nachrichten/pkg/http"
)

const (
	baseUrl      string = "https://www.tagesschau.de/"
	apiUrl       string = baseUrl + "api2u/homepage"
	shortNewsUrl string = baseUrl + "multimedia/sendung/tagesschau_in_100_sekunden"
)

type ImageSize int
type AspectRation int

const (
	SMALL ImageSize = iota
	MEDIUM
	LARGE

	SQUARE AspectRation = iota
	RECT
)

type ImageSpec struct {
	Size  ImageSize
	Ratio AspectRation
}

type News struct {
	NationalNews []Article `json:"news"`
	RegionalNews []Article `json:"regional"`
}

type Article struct {
	Topline      string     `json:"topline"`
	Desc         string     `json:"title"`
	Introduction string     `json:"firstSentence"`
	Tags         []Tag      `json:"tags"`
	Type         string     `json:"type"`
	Ressort      string     `json:"ressort"`
	RegionID     RegionID   `json:"regionId"`
	RegionIDs    []RegionID `json:"regionIds"`
	URL          string     `json:"shareURL"`
	Breaking     bool       `json:"breakingNews"`
	Date         time.Time  `json:"date"`
	ImageData    ImageData  `json:"teaserImage"`
	Content      []Content  `json:"content"`
	Video        Video      `json:"video"`
	ID           string     `json:"externalId"`
	Details      string     `json:"details"`
	DetailsWeb   string     `json:"detailsweb"`
}

func (n Article) Title() string       { return n.Topline }
func (n Article) Description() string { return n.Desc }
func (n Article) FilterValue() string { return n.Topline }

func (n Article) GetRelatedArticles() []Article {
	articles := []Article{}
	for _, content := range n.Content {
		if content.Type == "related" {
			articles = append(articles, content.Related...)
		}
	}
	return articles
}

func (n Article) IsRegionalArticle() bool {
	return len(n.RegionIDs) > 0
}

type Tag struct {
	Tag string `json:"tag"`
}

type Content struct {
	Value   string    `json:"value"`
	Type    string    `json:"type"`
	Related []Article `json:"related"`
}

type ImageData struct {
	Title         string        `json:"alttext"`
	Type          string        `json:"type"`
	ImageVariants ImageVariants `json:"imageVariants"`
}

type ImageVariants struct {
	SquareSmall  string `json:"1x1-144"`
	SquareMedium string `json:"1x1-432"`
	SquareLarge  string `json:"1x1-840"`
	RectSmall    string `json:"16x9-256"`
	RectMedium   string `json:"16x9-640"`
	RectLarge    string `json:"16x9-1920"`
}

type Video struct {
	Title         string        `json:"title"`
	Date          time.Time     `json:"date"`
	VideoVariants VideoVariants `json:"streams"`
}

type VideoVariants struct {
	Small  string `json:"h264s"`
	Medium string `json:"h264m"`
	Big    string `json:"h264xl"`
}

func RegionIdToName(id int) (string, error) {
	regionId := RegionID(id)
	regionName, ok := GERMAN_NAMES[regionId]
	if ok {
		return string(regionName), nil
	}
	return "", errors.New("invalid regionId")
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
	news.NationalNews = deduplicateArticles(news.NationalNews)
	news.RegionalNews = deduplicateArticles(news.RegionalNews)
	return news
}

func LoadArticle(url string) *Article {
	body, err := http.FetchURL(url)
	if err != nil {
		log.Fatal(err)
	}

	var article Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		log.Fatal(err)
	}
	return &article

}

func deduplicateArticles(articles []Article) []Article {
	deduped := []Article{}
	seen := make(map[string]bool)
	for _, entry := range articles {
		id := entry.ID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = true
		deduped = append(deduped, entry)
	}
	return deduped
}

func getImageURL(variants ImageVariants, imageSpec ImageSpec) string {
	sizeMap := map[ImageSpec]string{
		{SMALL, RECT}:    variants.RectSmall,
		{MEDIUM, RECT}:   variants.RectMedium,
		{LARGE, RECT}:    variants.RectLarge,
		{SMALL, SQUARE}:  variants.SquareSmall,
		{MEDIUM, SQUARE}: variants.SquareMedium,
		{LARGE, SQUARE}:  variants.SquareLarge,
	}
	return sizeMap[imageSpec]
}

func (news *News) GetArticlesOfRegion(regionId RegionID) []Article {
	allEntries := news.getCombinedArticles()
	entries := []Article{}
	for _, e := range allEntries {
		if contains(e.RegionIDs, regionId) {
			entries = append(entries, e)
		}
	}
	return entries
}

func (news *News) getCombinedArticles() []Article {
	entries := []Article{}
	entries = append(entries, news.NationalNews...)
	entries = append(entries, news.RegionalNews...)
	return entries
}

func contains(s []RegionID, e RegionID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
