package handlers

import (
	"net/http"

	"github.com/amanallah-jendoubi/Textio/http/helpers"
	"github.com/amanallah-jendoubi/Textio/http/middlewares"
	"github.com/amanallah-jendoubi/Textio/sql/database"
	"github.com/google/uuid"
)

func UserInfoHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, 500, "missing user id in context")
			return
		}
		userName, err := q.GetUserByID(r.Context(), userID)
		if err != nil {
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		type response struct {
			UserName string `json:"user_name"`
		}
		res := response{
			UserName: userName,
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}
