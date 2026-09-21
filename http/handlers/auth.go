package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/amanallah-jendoubi/Textio/http/helpers"
	"github.com/amanallah-jendoubi/Textio/sql/database"
	"github.com/google/uuid"
)

func RegistrationHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		type registrationRequest struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		var req registrationRequest
		decoder := json.NewDecoder(r.Body)
		// request validation
		if err := decoder.Decode(&req); err != nil {
			helpers.RespondWithError(w, 400, "invalid request body")
			return
		}
		if req.Name == "" || req.Password == "" {
			helpers.RespondWithError(w, 400, "name and password are required")
			return
		}
		//check if user already registred
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
		//create user
		user, err := q.CreateUser(r.Context(), database.CreateUserParams{Name: req.Name, Password: pwdHash})
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create user")
			return
		}
		//refresh token generation
		rawToken, hashedToken, err := helpers.GenerateRefreshToken()
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create refresh token")
			return
		}
		q.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: user.ID, Token: hashedToken, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), FamilyID: uuid.New()})

		//access token generation
		accessToken, err := helpers.GenerateAccessToken(user.ID)
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create access token")
			return
		}

		type response struct {
			ID           uuid.UUID `json:"id"`
			Name         string    `json:"name"`
			RefreshToken string    `json:"refresh_token"`
			AccessToken  string    `json:"access_token"`
		}
		res := response{
			ID:           user.ID,
			Name:         user.Name,
			RefreshToken: rawToken,
			AccessToken:  accessToken,
		}

		helpers.RespondWithJSON(w, 201, res)
	}
	return http.HandlerFunc(fn)
}

func LoginHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		type loginRequest struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		var req loginRequest
		decoder := json.NewDecoder(r.Body)
		// request validation
		if err := decoder.Decode(&req); err != nil {
			helpers.RespondWithError(w, 400, "invalid request body")
			return
		}
		if req.Name == "" || req.Password == "" {
			helpers.RespondWithError(w, 400, "name and password are required")
			return
		}
		//check if user already registred
		user, err := q.GetUserByName(r.Context(), req.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helpers.RespondWithError(w, 404, "verify you user name")
				return
			}
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		err = helpers.VerifPassword(user.Password, req.Password)
		if err != nil {
			helpers.RespondWithError(w, 401, "verify your password")
			return
		}
		// refresh token generation
		rawToken, hashedToken, err := helpers.GenerateRefreshToken()
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create refresh token")
			return
		}
		q.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: user.ID, Token: hashedToken, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), FamilyID: uuid.New()})

		//access token generation
		accessToken, err := helpers.GenerateAccessToken(user.ID)
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create access token")
			return
		}
		type response struct {
			ID           uuid.UUID `json:"id"`
			Name         string    `json:"name"`
			RefreshToken string    `json:"refresh_token"`
			AccessToken  string    `json:"access_token"`
		}
		res := response{
			ID:           user.ID,
			Name:         user.Name,
			RefreshToken: rawToken,
			AccessToken:  accessToken,
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}

func RefreshHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		type refreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}
		var req refreshRequest
		decoder := json.NewDecoder(r.Body)
		// request validation
		if err := decoder.Decode(&req); err != nil {
			helpers.RespondWithError(w, 400, "invalid request body")
			return
		}
		if req.RefreshToken == "" {
			helpers.RespondWithError(w, 400, "refresh token required")
			return
		}
		//refresh token db lookup
		refreshToken, err := q.GetRefreshToken(r.Context(), helpers.HashToken(req.RefreshToken))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helpers.RespondWithError(w, 404, "verify you user name")
				return
			}
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.UsedAt.Valid {
			err = q.RevokeRefreshTokensByFamilyID(r.Context(), refreshToken.FamilyID)
			if err != nil {
				helpers.RespondWithError(w, 500, "internal server error")
				return
			}
			helpers.RespondWithError(w, 401, "invalid or expired refresh token")
			return
		}
		//generate refresh token in the same family tree
		rawToken, hashedToken, err := helpers.GenerateRefreshToken()
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create refresh token")
			return
		}
		q.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: refreshToken.UserID, Token: hashedToken, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), FamilyID: refreshToken.FamilyID})
		//generate access token
		accessToken, err := helpers.GenerateAccessToken(refreshToken.UserID)
		if err != nil {
			helpers.RespondWithError(w, 500, "could not create access token")
			return
		}
		type response struct {
			RefreshToken string `json:"refresh_token"`
			AccessToken  string `json:"access_token"`
		}
		res := response{
			RefreshToken: rawToken,
			AccessToken:  accessToken,
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}
