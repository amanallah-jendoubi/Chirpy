package realtime

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu       sync.Mutex
	clients  map[uuid.UUID]*websocket.Conn
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[uuid.UUID]*websocket.Conn),
		upgrader: websocket.Upgrader{},
	}
}

func (h *Hub) Serve(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.mu.Lock()
	previous := h.clients[userID]
	h.clients[userID] = conn
	h.mu.Unlock()
	if previous != nil {
		_ = previous.Close()
	}

	defer func() {
		h.mu.Lock()
		if h.clients[userID] == conn {
			delete(h.clients, userID)
		}
		h.mu.Unlock()
		_ = conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) SendToUsers(userIDs []uuid.UUID, message any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, userID := range userIDs {
		conn, ok := h.clients[userID]
		if !ok {
			continue
		}
		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			delete(h.clients, userID)
			_ = conn.Close()
			continue
		}
		if err := conn.WriteJSON(message); err != nil {
			delete(h.clients, userID)
			_ = conn.Close()
		}
	}
}
