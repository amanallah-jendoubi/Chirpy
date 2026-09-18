package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/amanallah-jendoubi/Textio/db"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handlerLogger(mux),
	}
	rows, err := db.DB.Query("SELECT NOW()")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	var now string
	for rows.Next() {
		if err := rows.Scan(&now); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		log.Println("Current time from DB:", now)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration failed: %v", err)
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
