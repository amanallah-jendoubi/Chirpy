package auth

import (
	"encoding/json"
	"github.com/zalando/go-keyring"
	"io"
	"net/http"
)

const (
	service = "matrix_chat"
	user    = "current"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Save(t Tokens) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return keyring.Set(service, user, string(b))
}

func Load() (Tokens, error) {
	var t Tokens
	s, err := keyring.Get(service, user)
	if err != nil {
		return t, err // keyring.ErrNotFound => not logged in
	}
	return t, json.Unmarshal([]byte(s), &t)
}

func Clear() error {
	return keyring.Delete(service, user) // on logout
}

func StoreTokens(res *http.Response) error {
	var tokens Tokens
	decoder := json.NewDecoder(io.LimitReader(res.Body, 1024))
	if err := decoder.Decode(&tokens); err != nil {
		return err
	}
	if err := Save(tokens); err != nil {
		return err
	}
	return nil
}

func GetAccessToken() (string, error) {
	t, err := Load()
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}
