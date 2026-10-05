package client

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
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
