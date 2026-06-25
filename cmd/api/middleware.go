package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/0xdonnie7/Expense_API/internal/auth"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (app *application) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method:%v", t.Header["alg"])
			}
			return []byte(app.config.jwt.secret), nil
		})

		if err != nil || !token.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		userIDString, ok := claims["user_id"].(string)
		if !ok || userIDString == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		userID, err := uuid.Parse(userIDString)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		r = r.WithContext(auth.ContextWithUserID(r.Context(), userID))

		next.ServeHTTP(w, r)
	})
}
