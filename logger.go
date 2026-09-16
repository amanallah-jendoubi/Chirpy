package main

import (
	"log"
	"net/http"
)

func handlerLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
