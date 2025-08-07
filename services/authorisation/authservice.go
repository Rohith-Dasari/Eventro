package authorisation

import (
	"context"
	"errors"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo userrepository.UserRepository
}

func NewAuthService(userRepo userrepository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (a *AuthService) ValidateLogin(ctx context.Context, email, password string) (models.User, error) {
	for _, user := range a.UserRepo.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return models.User{}, errors.New("incorrect password")
			}
			return user, nil
		}
	}
	return models.User{}, errors.New("invalid email or password")
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (a *AuthService) UserExists(email string) bool {
	for _, user := range a.UserRepo.Users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (a *AuthService) IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
