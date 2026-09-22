package middlewares

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amanallah-jendoubi/Textio/http/helpers"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func VerifyAccessToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//get token_secret
		if err := godotenv.Load("./.env"); err != nil {
			log.Printf("failed to load .env: %v", err)
			helpers.RespondWithError(w, 500, "internal server error")
			return
		}
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		//get token from auth header
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			helpers.RespondWithError(w, 401, "missing authorization header")
			return
		}
		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}
		tokenString := parts[1]

		//parse access token
		var claims helpers.Claims
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			helpers.RespondWithError(w, 401, "invalid token")
			return
		}
		exp, err := token.Claims.GetExpirationTime()
		if err != nil || exp == nil || exp.Before(time.Now()) {
			helpers.RespondWithError(w, 401, "expired token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
