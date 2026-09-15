package main 

import (
	"fmt"
	"net/http"
)


func main {
	const port := "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(server.ListenAndServe())
}