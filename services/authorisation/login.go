package authorisation

import (
	"errors"
	"eventro/models"

	"golang.org/x/crypto/bcrypt"
)

// login
func ValidateLogin(users []models.User, email string, password string) (models.User, error) {
	for _, user := range users {
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

//hashing

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}
