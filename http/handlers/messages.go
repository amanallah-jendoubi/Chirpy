package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/amanallah-jendoubi/Textio/http/helpers"
	"github.com/amanallah-jendoubi/Textio/http/middlewares"
	"github.com/amanallah-jendoubi/Textio/sql/database"
)

func SendMessageHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		receiverIDString := r.PathValue("receiverID")
		if receiverIDString == "" {
			helpers.RespondWithError(w, 400, "recipientID is required")
			return
		}
		receiverID, err := uuid.Parse(receiverIDString)
		if err != nil {
			helpers.RespondWithError(w, 400, "invalid receiverID")
		}
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, 500, "missing user id in context")
			return
		}
		//if receiverID is chat_groupID then sender must be a member
		isChatGroupID, err := q.ChatGroupExists(r.Context(), receiverID)
		if err != nil {
			helpers.RespondWithError(w, 500, "server internal error")
			return
		}

		if isChatGroupID {
			isGroupMember, err := q.IsGroupMember(r.Context(), database.IsGroupMemberParams{UserID: userID, ChatGroupID: receiverID})
			if err != nil {
				helpers.RespondWithError(w, 500, "server internal error")
				return
			}
			if !isGroupMember {
				helpers.RespondWithError(w, 403, "not a group member")
				return
			}
		}
		type requset struct {
			Body string `json:"body"`
		}
		var req requset
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			helpers.RespondWithError(w, 400, "invalid request body")
			return
		}
		arg := database.CreateMessageParams{
			Body:       req.Body,
			SenderID:   userID,
			ReceiverID: receiverID,
		}
		message, err := q.CreateMessage(r.Context(), arg)
		if err != nil {
			helpers.RespondWithError(w, 500, "server internal error")
			return
		}
		type response struct {
			ID         uuid.UUID `json:"id"`
			CreatedAt  time.Time `json:"created_at"`
			SenderID   uuid.UUID `json:"sender_id"`
			ReceiverID uuid.UUID `json:"receiver_id"`
			Body       string    `json:"body"`
		}
		helpers.RespondWithJSON(w, 200, response{
			ID:         message.ID,
			CreatedAt:  message.CreatedAt,
			SenderID:   message.SenderID,
			ReceiverID: message.ReceiverID,
			Body:       message.Body,
		})
	}
	return http.HandlerFunc(fn)
}

func ConversationsHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, 500, "missing user id in context")
			return
		}
		conversations, err := q.GetConversationsByUserID(r.Context(), userID)
		if err != nil {
			helpers.RespondWithError(w, 500, "server internal error")
			return
		}
		type response struct {
			Name   string      `json:"name"`
			Latest interface{} `json:"latest"`
		}
		var res []response
		for _, c := range conversations {
			res = append(res, response{Name: c.Name, Latest: c.Latest})
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}

func GetMessagesHandler(q *database.Queries) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		receiverIDString := r.PathValue("receiverID")
		if receiverIDString == "" {
			helpers.RespondWithError(w, 400, "recipientID is required")
			return
		}
		receiverID, err := uuid.Parse(receiverIDString)
		if err != nil {
			helpers.RespondWithError(w, 400, "invalid receiverID")
		}
		userID, ok := r.Context().Value(middlewares.UserIDContextKey).(uuid.UUID)
		if !ok {
			helpers.RespondWithError(w, 500, "missing user id in context")
			return
		}
		//if receiverID is chat_groupID then sender must be a member
		isChatGroupID, err := q.ChatGroupExists(r.Context(), receiverID)
		if err != nil {
			helpers.RespondWithError(w, 500, "server internal error")
			return
		}
		type response struct {
			ID         uuid.UUID `json:"id"`
			CreatedAt  time.Time `json:"created_at"`
			SenderID   uuid.UUID `json:"sender_id"`
			ReceiverID uuid.UUID `json:"receiver_id"`
			Body       string    `json:"body"`
		}
		var res []response
		if isChatGroupID {
			isGroupMember, err := q.IsGroupMember(r.Context(), database.IsGroupMemberParams{UserID: userID, ChatGroupID: receiverID})
			if err != nil {
				helpers.RespondWithError(w, 500, "server internal error")
				return
			}
			if !isGroupMember {
				helpers.RespondWithError(w, 403, "not a group member")
				return
			}
			messages, err := q.GetGroupMessages(r.Context(), receiverID)
			if err != nil {
				helpers.RespondWithError(w, 500, "server internal error")
				return
			}
			for _, m := range messages {
				res = append(res, response{
					ID:         m.ID,
					CreatedAt:  m.CreatedAt,
					SenderID:   m.SenderID,
					ReceiverID: m.ReceiverID,
					Body:       m.Body,
				})
			}
			helpers.RespondWithJSON(w, 200, res)
			return
		}
		messages, err := q.GetDuelMessages(r.Context(), database.GetDuelMessagesParams{ReceiverID: receiverID, SenderID: userID})
		if err != nil {
			helpers.RespondWithError(w, 500, "server internal error")
			return
		}
		for _, m := range messages {
			res = append(res, response{
				ID:         m.ID,
				CreatedAt:  m.CreatedAt,
				SenderID:   m.SenderID,
				ReceiverID: m.ReceiverID,
				Body:       m.Body,
			})
		}
		helpers.RespondWithJSON(w, 200, res)
	}
	return http.HandlerFunc(fn)
}
