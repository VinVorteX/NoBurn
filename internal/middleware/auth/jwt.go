package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/VinVorteX/NoBurn/internal/utils"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(w, http.StatusUnauthorized, "Missing or invalid token")
			return
		}
		
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.VerifyToken(tokenStr)

		if err != nil || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}