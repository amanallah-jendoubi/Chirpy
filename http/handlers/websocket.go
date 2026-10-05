package handlers

import (
	"net/http"

	"github.com/amanallah-jendoubi/matrix-chat/http/middlewares"
	"github.com/amanallah-jendoubi/matrix-chat/http/realtime"
	"github.com/google/uuid"
)

func WebSocketHandler(hub *realtime.Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			http.Error(w, "missing user id in context", http.StatusInternalServerError)
			return
		}
		hub.Serve(w, r, userID)
	})
}
