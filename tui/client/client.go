package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/amanallah-jendoubi/Textio/tui/auth"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const ApiURL = "http://localhost:8080/api"

func NewClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

type tokenStore interface {
	Load() (auth.Tokens, error)
	Save(auth.Tokens) error
	Clear() error
}

type keyringTokenStore struct{}

func (keyringTokenStore) Load() (auth.Tokens, error)    { return auth.Load() }
func (keyringTokenStore) Save(tokens auth.Tokens) error { return auth.Save(tokens) }
func (keyringTokenStore) Clear() error                  { return auth.Clear() }

const sessionExpiredHeader = "X-Matrix-Session-Expired"

var refreshMu sync.Mutex

func Post(url string, accessToken string, body io.Reader) (*http.Response, error) {
	return requestWithRefresh(NewClient(), keyringTokenStore{}, ApiURL, http.MethodPost, url, accessToken, body)
}

func Get(url string, accessToken string, body io.Reader) (*http.Response, error) {
	return requestWithRefresh(NewClient(), keyringTokenStore{}, ApiURL, http.MethodGet, url, accessToken, body)
}

func requestWithRefresh(httpClient *http.Client, tokens tokenStore, baseURL, method, path, accessToken string, body io.Reader) (*http.Response, error) {
	var requestBody []byte
	var err error
	if body != nil {
		requestBody, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}
	response, err := doRequest(httpClient, baseURL, method, path, accessToken, requestBody)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusUnauthorized || accessToken == "" || path == "/refresh" {
		return response, nil
	}

	newAccessToken, sessionExpired, err := refreshAccessToken(httpClient, tokens, baseURL, accessToken)
	if err != nil {
		if sessionExpired {
			response.Header.Set(sessionExpiredHeader, "true")
		}
		return response, nil
	}
	_ = response.Body.Close()
	return doRequest(httpClient, baseURL, method, path, newAccessToken, requestBody)
}

func doRequest(httpClient *http.Client, baseURL, method, path, accessToken string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return httpClient.Do(req)
}

func refreshAccessToken(httpClient *http.Client, tokens tokenStore, baseURL, rejectedToken string) (string, bool, error) {
	refreshMu.Lock()
	defer refreshMu.Unlock()

	stored, err := tokens.Load()
	if err != nil {
		return "", false, err
	}
	if stored.AccessToken != "" && stored.AccessToken != rejectedToken {
		return stored.AccessToken, false, nil
	}
	if stored.RefreshToken == "" {
		return "", false, fmt.Errorf("refresh token is missing")
	}

	body, err := json.Marshal(struct {
		RefreshToken string `json:"refresh_token"`
	}{RefreshToken: stored.RefreshToken})
	if err != nil {
		return "", false, err
	}
	response, err := doRequest(httpClient, baseURL, http.MethodPost, "/refresh", "", body)
	if err != nil {
		return "", false, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		sessionExpired := response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusNotFound
		if sessionExpired {
			_ = tokens.Clear()
		}
		return "", sessionExpired, fmt.Errorf("refresh endpoint returned HTTP %d", response.StatusCode)
	}

	var refreshed auth.Tokens
	if err := json.NewDecoder(io.LimitReader(response.Body, 1024)).Decode(&refreshed); err != nil {
		return "", false, err
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		return "", false, fmt.Errorf("refresh response is missing tokens")
	}
	if err := tokens.Save(refreshed); err != nil {
		return "", false, err
	}
	return refreshed.AccessToken, false, nil
}

type WebSocket struct {
	conn *websocket.Conn
}

type RealtimeMessage struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Body       string    `json:"body"`
}

func ConnectWebSocket(accessToken string) (*WebSocket, error) {
	endpoint, err := url.Parse(ApiURL)
	if err != nil {
		return nil, err
	}
	if endpoint.Scheme == "https" {
		endpoint.Scheme = "wss"
	} else {
		endpoint.Scheme = "ws"
	}
	endpoint.Path = "/ws"

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 5 * time.Second
	conn, _, err := dialer.Dial(endpoint.String(), http.Header{
		"Authorization": []string{"Bearer " + accessToken},
	})
	if err != nil {
		return nil, err
	}
	return &WebSocket{conn: conn}, nil
}

func (ws *WebSocket) ReadMessage() (RealtimeMessage, error) {
	var message RealtimeMessage
	err := ws.conn.ReadJSON(&message)
	return message, err
}

func (ws *WebSocket) Close() error {
	return ws.conn.Close()
}
