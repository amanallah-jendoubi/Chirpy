package client

import (
	"io"
	"net/http"
	"time"
)

const ApiURL = "http://localhost:8080/api"

func NewClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func Post(url string, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, ApiURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return NewClient().Do(req)
}

func Get(url string, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", ApiURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return NewClient().Do(req)
}
