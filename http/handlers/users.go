package handlers

import (
	"database/sql"
	"net/http"

	"github.com/amanallah-jendoubi/matrix-chat/http/helpers"
	"github.com/amanallah-jendoubi/matrix-chat/http/middlewares"
	"github.com/amanallah-jendoubi/matrix-chat/sql/database"
	"github.com/google/uuid"
)

func UserInfoHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, 500, "missing user id in context")
			return
		}
		user, err := q.GetUserByID(r.Context(), userID)
		if err != nil {
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		type response struct {
			UserID   uuid.UUID `json:"user_id"`
			UserName string    `json:"user_name"`
		}
		res := response{
			UserName: user.Name,
			UserID:   user.ID,
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}

func GetUserByIDHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.PathValue("userID"))
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, "invalid userID")
			return
		}
		user, err := q.GetUserByID(r.Context(), userID)
		if err != nil {
			if err == sql.ErrNoRows {
				helpers.RespondWithError(w, http.StatusNotFound, "user not found")
				return
			}
			helpers.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		type response struct {
			UserID   uuid.UUID `json:"user_id"`
			UserName string    `json:"user_name"`
		}
		helpers.RespondWithJSON(w, http.StatusOK, response{UserID: user.ID, UserName: user.Name})
	}
	return http.HandlerFunc(fn)
}

func GetUserIdByUserNameHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userName := r.PathValue("userName")
		userID, err := q.GetUserIDByUserName(r.Context(), userName)
		if err != nil {
			if err == sql.ErrNoRows {
				helpers.RespondWithError(w, http.StatusNotFound, "user not found")
				return
			}
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		type response struct {
			UserID uuid.UUID `json:"user_id"`
		}
		res := response{
			UserID: userID,
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}
