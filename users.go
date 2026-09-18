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
		helpers.RespondWithError(w, 400, "Invalid user creation request")
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
	helpers.RespondWithJSON(w, 201, helpers.DatabaseUserToUser(user))
}
