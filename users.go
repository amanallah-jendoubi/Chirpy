package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	dbgen "github.com/amanallah-jendoubi/Textio/db/sqlc"
	"github.com/amanallah-jendoubi/Textio/helpers"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := createUserRequest{}
	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), dbgen.CreateUserParams{
		Email: params.Email,
		Name:  sql.NullString{String: params.Name, Valid: true},
	})
	if err != nil {
		helpers.RespondWithError(w, 500, "could not create user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "could not encode user", http.StatusInternalServerError)
		return
	}
}
