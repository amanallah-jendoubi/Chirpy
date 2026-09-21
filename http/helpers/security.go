package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GenerateRefreshToken returns:
//   - rawToken: the base64url-encoded 256-bit random value to send to the client
//   - hashedToken: the SHA-256 hash (hex-encoded) to store via CreateRefreshToken
func GenerateRefreshToken() (rawToken string, hashedToken string, err error) {
	buf := make([]byte, 32) // 256 bits
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}

	rawToken = base64.RawURLEncoding.EncodeToString(buf)

	sum := sha256.Sum256([]byte(rawToken))
	hashedToken = hex.EncodeToString(sum[:])

	return rawToken, hashedToken, nil
}
