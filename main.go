package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "./public"
	const port = "8080"

	cfg := &apiConfig{}

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/healthz/", handlerReadiness)
	mux.Handle("/metrics/", cfg.handlerMetrics())
	mux.Handle("/reset/", cfg.handlerReset())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLog(mux),
	}

	log.Printf("Running on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

//handlers

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, "Hits: %d", cfg.fileserverHits.Load()))
	}
}

func (cfg *apiConfig) handlerReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
	}
}

//middlewares

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
