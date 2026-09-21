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

	log.Print("Listening...")
	http.ListenAndServe(":8080", middlewares.Logger(mux))
}
