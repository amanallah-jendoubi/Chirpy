package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/amanallah-jendoubi/Textio/http/helpers"
	"github.com/amanallah-jendoubi/Textio/http/middlewares"
	"github.com/amanallah-jendoubi/Textio/sql/database"
	"github.com/google/uuid"
)

func CreateGroupHandler(q *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, http.StatusInternalServerError, "missing user id in context")
			return
		}

		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		name := strings.TrimSpace(req.Name)
		if name == "" || len([]rune(name)) > 80 {
			helpers.RespondWithError(w, http.StatusBadRequest, "group name must be between 1 and 80 characters")
			return
		}

		group, err := q.CreateChatGroupWithCreator(r.Context(), database.CreateChatGroupWithCreatorParams{
			Name:      name,
			CreatedBy: userID,
		})
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "could not create group")
			return
		}

		helpers.RespondWithJSON(w, http.StatusCreated, struct {
			ID        uuid.UUID `json:"id"`
			Name      string    `json:"name"`
			CreatedBy uuid.UUID `json:"created_by"`
		}{ID: group.ID, Name: group.Name, CreatedBy: group.CreatedBy})
	})
}

func AddGroupMemberHandler(q *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, http.StatusInternalServerError, "missing user id in context")
			return
		}

		groupID, err := uuid.Parse(r.PathValue("groupID"))
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, "invalid group id")
			return
		}
		var req struct {
			UserID uuid.UUID `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == uuid.Nil {
			helpers.RespondWithError(w, http.StatusBadRequest, "valid user_id is required")
			return
		}

		exists, err := q.ChatGroupExists(r.Context(), groupID)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "could not check group")
			return
		}
		if !exists {
			helpers.RespondWithError(w, http.StatusNotFound, "group not found")
			return
		}

		isMember, err := q.IsGroupMember(r.Context(), database.IsGroupMemberParams{
			UserID:      userID,
			ChatGroupID: groupID,
		})
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "could not check group membership")
			return
		}
		if !isMember {
			helpers.RespondWithError(w, http.StatusForbidden, "only group members can add members")
			return
		}

		if _, err := q.GetUserByID(r.Context(), req.UserID); err != nil {
			if err == sql.ErrNoRows {
				helpers.RespondWithError(w, http.StatusNotFound, "user not found")
				return
			}
			helpers.RespondWithError(w, http.StatusInternalServerError, "could not check user")
			return
		}

		if err := q.AddChatGroupMember(r.Context(), database.AddChatGroupMemberParams{
			UserID:      req.UserID,
			ChatGroupID: groupID,
		}); err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "could not add group member")
			return
		}

		helpers.RespondWithJSON(w, http.StatusOK, struct {
			GroupID uuid.UUID `json:"group_id"`
			UserID  uuid.UUID `json:"user_id"`
		}{GroupID: groupID, UserID: req.UserID})
	})
}
