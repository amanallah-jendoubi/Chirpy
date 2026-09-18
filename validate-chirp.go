package main

import (
	"encoding/json"
	"github.com/amanallah-jendoubi/Textio/helpers"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	dec := json.NewDecoder(r.Body)
	var params chirpRequest
	if err := dec.Decode(&params); err != nil {
		helpers.RespondWithError(w, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		helpers.RespondWithError(w, 400, "Chirp is too long !")
		return
	}
	helpers.RespondWithJSON(w, 200, response{CleanedBody: cleanBody(params.Body)})
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
