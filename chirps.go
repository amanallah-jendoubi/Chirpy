package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	dbgen "github.com/amanallah-jendoubi/Textio/db/sqlc"
	"github.com/amanallah-jendoubi/Textio/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := createChirpRequest{}
	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("error in chirp request creation %v", err))
		return
	} else if len(params.Body) > 140 {
		helpers.RespondWithError(w, 400, "Chirp is too long !")
		return
	}
	chirp, err := cfg.dbQueries.CreateChirps(r.Context(), dbgen.CreateChirpsParams{
		Body:   cleanBody(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("could not create chirp %v", err))
		return
	}
	helpers.RespondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("could not get chirps %v", err))
		return
	}

	helpers.RespondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		helpers.RespondWithError(w, 400, "invalid chirp id")
		return
	}

	chirp, err := cfg.dbQueries.GetChirpById(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.RespondWithError(w, 404, "chirp not found")
			return
		}
		helpers.RespondWithError(w, 500, fmt.Sprintf("%v", err))
		return
	}

	helpers.RespondWithJSON(w, 200, chirp)
}

func cleanBody(body string) (cleanedBody string) {
	s := strings.ToLower(body)
	words := strings.Split(s, " ")
	for i, w := range words {
		if w == "kerfuffle" || w == "sharbert" || w == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
