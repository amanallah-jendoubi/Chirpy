package main

import (
	"encoding/json"
	"github.com/amanallah-jendoubi/Textio/helpers"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
	}
	dec := json.NewDecoder(r.Body)
	var params chirpRequest
	err := dec.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		helpers.RespondWithError(w, 400, "Chirp is too long !")
		return
	}
	helpers.RespondWithJSON(w, 200, validResponse{Valid: true})
}
