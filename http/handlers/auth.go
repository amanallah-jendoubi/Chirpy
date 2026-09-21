package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amanallah-jendoubi/Textio/http/helpers"
	"github.com/amanallah-jendoubi/Textio/sql/database"
	"github.com/google/uuid"
	"time"
)

func RegistrationHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		type registrationRequest struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		var req registrationRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			helpers.RespondWithError(w, 400, "invalid request body")
			return
		}
		if req.Name == "" || req.Password == "" {
			helpers.RespondWithError(w, 400, "name and password are required")
			return
		}
		exists, err := q.UserExists(r.Context(), req.Name)
		if err != nil {
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		if exists {
			helpers.RespondWithError(w, 409, "name already used")
			return
		}
		pwdHash, err := helpers.HashPassword(req.Password)
		if err != nil {
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		user, err := q.CreateUser(r.Context(), database.CreateUserParams{Name: req.Name, Password: pwdHash})
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create user")
			return
		}
		rawToken, hashedToken, err := helpers.GenerateRefreshToken()
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create refresh token")
			return
		}
		q.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: user.ID, Token: hashedToken, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), FamilyID: uuid.New()})
		type response struct {
			ID           uuid.UUID `json:"id"`
			Name         string    `json:"name"`
			RefreshToken string    `json:"refresh_token"`
		}
		res := response{
			ID:           user.ID,
			Name:         user.Name,
			RefreshToken: rawToken,
		}

		//to do access token

		helpers.RespondWithJSON(w, 201, res)
	}
	return http.HandlerFunc(fn)
}
