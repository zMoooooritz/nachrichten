package http

import (
	"image"
	"io"
	"net/http"
	"time"
)

const (
	httpTimeout time.Duration = 2
	agentName   string        = "nachrichten-agent"
)

var (
	client http.Client = http.Client{
		Timeout: time.Second * httpTimeout,
	}
)

func FetchURL(url string) ([]byte, error) {
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

func LoadImage(url string) (image.Image, error) {
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

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
