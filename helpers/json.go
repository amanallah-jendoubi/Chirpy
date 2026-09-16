package helpers

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	dat, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	type errorResponse struct {
		Error string `json:"error"`
	}
	return RespondWithJSON(w, code, errorResponse{Error: msg})
}
