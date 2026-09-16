package main

import (
	"fmt"
	"net/http"
)

func handlerLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, "request: %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
