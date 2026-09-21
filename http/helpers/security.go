package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// password hadsh
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

//verify password

func VerifPassword(hashedPwd, rawPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(rawPwd))
	if err != nil {
		return err
	}
	return nil
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

// access token generation

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID) (string, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	now := time.Now()
	claims := Claims{
		UserID:           userID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)), IssuedAt: jwt.NewNumericDate(now)},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
