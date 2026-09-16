package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const filepathRoot = "./public"
	const port = "8080"

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.Handle("/time", timeHandler(time.RFC1123))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLog(mux),
	}

	log.Printf("Running on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func timeHandler(format string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("the time is: " + time.Now().Format(format)))
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
