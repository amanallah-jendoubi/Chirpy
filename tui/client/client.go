package client

import (
	"io"
	"net/http"
	"time"
)

type ErrMsg struct{ Err error }

func (e ErrMsg) Error() string { return e.Err.Error() }

const ApiURL = "http://localhost:8080/api"

func NewClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func Post(url string, body io.Reader) (*http.Response, error) {
	return NewClient().Post(ApiURL+url, "application/json", body)
}
