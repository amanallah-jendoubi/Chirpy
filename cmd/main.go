package main

import (
	"github.com/amanallah-jendoubi/Textio/config"
	"github.com/amanallah-jendoubi/Textio/http/handlers"
	"github.com/amanallah-jendoubi/Textio/http/middlewares"
	"github.com/amanallah-jendoubi/Textio/sql/database"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()

	envPath := filepath.Join(".", ".env")
	config.Init(envPath)

	pool := config.DB
	queries := database.New(pool)

	//auth
	mux.Handle("POST /api/register", handlers.RegistrationHandler(queries))
	mux.Handle("POST /api/login", handlers.LoginHandler(queries))
	mux.Handle("POST /api/refresh", handlers.RefreshHandler(queries))
	//get current user info
	mux.Handle("GET /api/users/me", middlewares.VerifyAccessToken(handlers.UserInfoHandler(queries)))
	//send a message to user or chat group
	mux.Handle("POST /api/{receiverID}/messages", middlewares.VerifyAccessToken(handlers.SendMessageHandler(queries)))
	// get all conversations sorted by latest
	mux.Handle("GET /api/conversations", middlewares.VerifyAccessToken(handlers.ConversationsHandler(queries)))
	// get conversation messages (to improve)
	mux.Handle("GET /api/{receiverID}/messages", middlewares.VerifyAccessToken(handlers.GetMessagesHandler(queries)))
	/*todo
	create group
	add member to a group
	*/
	log.Print("Listening...")
	http.ListenAndServe(":8080", middlewares.Logger(mux))
}
