package http

import (
	"io"
	"net/http"
	"time"
)

const (
	httpTimeout time.Duration = 2
	agentName   string        = "nachrichten-agent"
)

func FetchURL(url string) ([]byte, error) {
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
