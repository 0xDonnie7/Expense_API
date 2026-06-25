package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type ContextKey string

const userIDKey ContextKey = "user_id"

func GenerateJWT(userID uuid.UUID, email, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthenticatedUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	return userID, ok
}

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
