package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type errBase struct {
	UserErr string // shown in the UI
	LogErr  error  // goes to logs
}

func (e errBase) Error() string { return e.UserErr }

type ListErrMsg struct{ errBase }
type ChatErrMsg struct{ errBase }
type GroupErrMsg struct{ errBase }
type DMErrMsg struct{ errBase }
type LoginErrMsg struct{ errBase }
type SignupErrMsg struct{ errBase }
type AuthErrMsg struct{ errBase }
type SessionExpiredMsg struct{}

// constructors, one per target
func AsList(b errBase) tea.Msg   { return ListErrMsg{b} }
func AsChat(b errBase) tea.Msg   { return ChatErrMsg{b} }
func AsGroup(b errBase) tea.Msg  { return GroupErrMsg{b} }
func AsDM(b errBase) tea.Msg     { return DMErrMsg{b} }
func AsAuth(b errBase) tea.Msg   { return AuthErrMsg{b} }
func AsLogin(b errBase) tea.Msg  { return LoginErrMsg{b} }
func AsSignup(b errBase) tea.Msg { return SignupErrMsg{b} }

// failure on frontend side (unreachable, bad JSON...)
func HandleLocalError(wrap func(errBase) tea.Msg, logErr error) tea.Msg {
	log.Println("local error:", logErr)
	return wrap(errBase{UserErr: "Something went wrong. Please try again.", LogErr: logErr})
}

// backend answered with a non-2xx status
func statusMessage(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "Something looks wrong with your input."
	case http.StatusUnauthorized:
		return "Your session expired. Please log in again."
	case http.StatusForbidden:
		return "You don't have access to this."
	case http.StatusNotFound:
		return "We couldn't find what you were looking for."
	case http.StatusConflict:
		return "That already exists."
	case http.StatusTooManyRequests:
		return "Too many requests. Please wait a moment."
	}
	if code >= 500 {
		return "The server had a problem. Please try again in a moment."
	}
	return "Something went wrong. Please try again."
}

func HandleAPIError(wrap func(errBase) tea.Msg, res *http.Response, fallback string) tea.Msg {
	defer res.Body.Close()
	if res.Header.Get("X-Matrix-Session-Expired") == "true" {
		return SessionExpiredMsg{}
	}

	var payload struct {
		Error string `json:"error"`
	}
	// decode failure is fine: empty or non-JSON body falls through to the defaults
	_ = json.NewDecoder(io.LimitReader(res.Body, 1024)).Decode(&payload)

	backendMsg := strings.TrimSpace(payload.Error)

	var text string
	switch {
	case fallback != "": // caller knows the context best
		text = fallback
	case backendMsg == "" || res.StatusCode >= 500: // never show raw 5xx text
		text = statusMessage(res.StatusCode)
	default: // 4xx with a message from the backend
		text = backendMsg
	}

	logErr := fmt.Errorf("http %d %s %s: %q", res.StatusCode, res.Request.Method, res.Request.URL, payload.Error)
	return wrap(errBase{UserErr: text, LogErr: logErr})
}
