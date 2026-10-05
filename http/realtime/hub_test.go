package realtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func TestSendToUsersDeliversToConnectedUser(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Serve(w, r, userID)
	}))
	defer server.Close()

	endpoint := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		t.Fatalf("connect WebSocket: %v", err)
	}
	defer conn.Close()

	deadline := time.Now().Add(time.Second)
	for {
		hub.mu.Lock()
		_, connected := hub.clients[userID]
		hub.mu.Unlock()
		if connected {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("user was not registered in hub")
		}
		time.Sleep(time.Millisecond)
	}

	hub.SendToUsers([]uuid.UUID{userID}, map[string]string{"body": "hello"})
	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("set read deadline: %v", err)
	}
	var message struct {
		Body string `json:"body"`
	}
	if err := conn.ReadJSON(&message); err != nil {
		t.Fatalf("read delivered message: %v", err)
	}
	if message.Body != "hello" {
		t.Fatalf("got body %q, want %q", message.Body, "hello")
	}
}
