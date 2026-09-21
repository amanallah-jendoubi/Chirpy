package middlewares

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URI: %s, Method: %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
