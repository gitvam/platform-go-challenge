package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gitvam/platform-go-challenge/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "my_super_secret"

type contextKey string

const contextKeyUserID = contextKey("userID")

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteJSONError(w, "unauthorized: invalid or missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		log.Println("[DEBUG] Raw token string:", tokenStr)
		log.Println("[DEBUG] Secret being used:", jwtSecret)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		log.Printf("[DEBUG] JWT parse result: valid=%v, err=%v\n", token.Valid, err)
		if err != nil || !token.Valid {
			utils.WriteJSONError(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["sub"] == nil {
			utils.WriteJSONError(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok {
			utils.WriteJSONError(w, "userID claim missing", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves the userID set by the JWT middleware
func GetUserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(contextKeyUserID).(string)
	return userID, ok
}
