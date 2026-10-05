package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/amanallah-jendoubi/matrix-chat/config"
	"github.com/amanallah-jendoubi/matrix-chat/http/handlers"
	"github.com/amanallah-jendoubi/matrix-chat/http/middlewares"
	"github.com/amanallah-jendoubi/matrix-chat/http/realtime"
	"github.com/amanallah-jendoubi/matrix-chat/sql/database"
)

func main() {
	mux := http.NewServeMux()

	envPath := filepath.Join(".", ".env")
	config.Init(envPath)

	pool := config.DB
	queries := database.New(pool)
	hub := realtime.NewHub()

	//auth
	mux.Handle("POST /api/register", handlers.RegistrationHandler(queries))
	mux.Handle("POST /api/login", handlers.LoginHandler(queries))
	mux.Handle("POST /api/refresh", handlers.RefreshHandler(queries))
	//get current user info
	mux.Handle("GET /api/users/me", middlewares.VerifyAccessToken(handlers.UserInfoHandler(queries)))
	mux.Handle("GET /api/users/id/{userID}", middlewares.VerifyAccessToken(handlers.GetUserByIDHandler(queries)))
	// get user ID by user name
	mux.Handle("GET /api/users/{userName}", middlewares.VerifyAccessToken(handlers.GetUserIdByUserNameHandler(queries)))
	// websocket upgrade
	mux.Handle("GET /ws", middlewares.VerifyAccessToken(handlers.WebSocketHandler(hub)))
	//send a message to user or chat group
	mux.Handle("POST /api/conversations/{receiverID}/messages", middlewares.VerifyAccessToken(handlers.SendMessageHandler(queries, hub)))
	// get all conversations sorted by latest
	mux.Handle("GET /api/conversations", middlewares.VerifyAccessToken(handlers.ConversationsHandler(queries)))
	// get conversation messages sorted by latest (to improve)
	mux.Handle("GET /api/conversations/{receiverID}/messages", middlewares.VerifyAccessToken(handlers.GetMessagesHandler(queries)))
	// group management
	mux.Handle("POST /api/groups", middlewares.VerifyAccessToken(handlers.CreateGroupHandler(queries)))
	mux.Handle("POST /api/groups/{groupID}/members", middlewares.VerifyAccessToken(handlers.AddGroupMemberHandler(queries)))
	log.Print("Listening...")
	http.ListenAndServe(":8080", middlewares.Logger(mux))
}
