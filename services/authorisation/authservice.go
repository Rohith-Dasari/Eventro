package authorisation

import (
	"context"
	"errors"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"net/mail"
	"regexp"

	"github.com/google/uuid"
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
	user, err := a.UserRepo.GetByEmail(email)
	if err != nil {
		return models.User{}, errors.New("invalid email or password")
	}
	if user.IsBlocked {
		return models.User{}, errors.New("user account is blocked, please contact admin")
	}

	// Compare entered password with stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, errors.New("invalid email or password")
	}

	return *user, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (a *AuthService) IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (a *AuthService) IsValidPhoneNumber(phone string) bool {
	var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	return phoneRegex.MatchString(phone)
}

func (a *AuthService) IsValidPassword(password string) bool {
	var (
		minLength = len(password) >= 12
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower  = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasSymbol = regexp.MustCompile(`[!@#$%^&*()\-+]`).MatchString(password)
	)

	return minLength && hasNumber && hasUpper && hasLower && hasSymbol
}

func (a *AuthService) Signup(ctx context.Context, username, email, phoneNumber, password string) (models.User, error) {
	if !a.IsValidEmail(email) {
		return models.User{}, errors.New("invalid email format")
	}

	if _, err := a.UserRepo.GetByEmail(email); err == nil {
		return models.User{}, errors.New("email already in use")
	}

	if !a.IsValidPassword(password) {
		return models.User{}, errors.New("password must be at least 12 characters long, and include uppercase, lowercase, number, and symbol")
	}

	hashedPassword, err := a.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	if !a.IsValidPhoneNumber(phoneNumber) {
		return models.User{}, errors.New("invalid phone number format")
	}

	newUser := models.User{
		UserID:      uuid.New().String(),
		Username:    username,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    hashedPassword,
		Role:        models.Customer,
	}

	if err := a.UserRepo.Create(&newUser); err != nil {
		return models.User{}, err
	}

	return newUser, nil
}
