package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/amanallah-jendoubi/Textio/db"
	dbgen "github.com/amanallah-jendoubi/Textio/db/sqlc"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *dbgen.Queries
}

func main() {
	db.Init(".env")
	const filepathRoot = "./public"
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/user", handlerCreateUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handlerLogger(mux),
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
