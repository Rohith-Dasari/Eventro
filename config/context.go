package config

import (
	"context"
	"eventro/models"
)

type contextKey string

const (
	UserIDKey   contextKey = "UserID"
	UserRoleKey contextKey = "Role"
	UserMailID  contextKey = "Email"
)

func GetUserID(ctx context.Context) string {
	id, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func GetUserRole(ctx context.Context) models.Role {
	role, ok := ctx.Value(UserRoleKey).(models.Role)
	if !ok {
		return ""
	}
	return role
}
