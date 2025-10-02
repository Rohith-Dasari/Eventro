package middleware

import (
	"context"
	"errors"
	"eventro2/services/authorisation"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserIDKey    contextKey = "userID"
	ContextUserEmailKey contextKey = "userEmail"
	ContextUserRoleKey  contextKey = "userRole"
)

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := authorisation.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, ContextUserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, ContextUserRoleKey, claims.Role)

		next(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ContextUserIDKey).(string)
	if !ok {
		return "", errors.New("userID not found in context")
	}
	return userID, nil
}

func GetUserEmail(ctx context.Context) (string, error) {
	email, ok := ctx.Value(ContextUserEmailKey).(string)
	if !ok {
		return "", errors.New("userEmail not found in context")
	}
	return email, nil
}

func GetUserRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value(ContextUserRoleKey).(string)
	if !ok {
		return "", errors.New("userRole not found in context")
	}
	return role, nil
}
