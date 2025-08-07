package authorisation

import (
	"context"
	"eventro2/models"
)

type AuthServiceInterface interface {
	ValidateLogin(ctx context.Context, email, password string) (models.User, error)
	HashPassword(password string) (string, error)
	UserExists(email string) bool
	IsValidEmail(email string) bool
}
