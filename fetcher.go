package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	baseUrl      string        = "https://www.tagesschau.de/"
	apiUrl       string        = baseUrl + "api2u/homepage"
	shortNewsUrl string        = baseUrl + "multimedia/sendung/tagesschau_in_100_sekunden"
	httpTimeout  time.Duration = 2
	agentName    string        = "Nachrichten-Agent"
)

type loadedNews News

func getNews() tea.Cmd {
	return func() tea.Msg {
		body, err := fetchURL(apiUrl)
		if err != nil {
			log.Fatal(err)
		}

		var news News
		err = json.Unmarshal(body, &news)
		if err != nil {
			log.Fatal(err)
		}
		return loadedNews(news)
	}
}

func getShortNewsURL() (string, error) {
	body, err := fetchURL(shortNewsUrl)
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

func fetchURL(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * httpTimeout,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", agentName)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
