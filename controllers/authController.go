package controllers

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"eventro2/services/authorisation"
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"golang.org/x/term"
)

func LoginFlow(ctx context.Context) context.Context {
	userRepo := userrepository.NewUserRepository()
	authService := authorisation.NewAuthService(userRepo)

	var user models.User

	var email, password string
	for {

		fmt.Print("Enter email: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return ctx
		}
		email = strings.TrimSpace(input)

		fmt.Print("Enter password: ")
		// input, err = reader.ReadString('\n')
		// if err != nil {
		// 	fmt.Println("Error reading input:", err)
		// 	return ctx
		// }
		//password = strings.TrimSpace(input)

		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\nError reading password:", err)
			return ctx
		}
		password = string(bytePassword)

		user, err = authService.ValidateLogin(ctx, email, password)
		if !user.IsBlocked {
			ctx = context.WithValue(ctx, config.UserIDKey, user.UserID)
			ctx = context.WithValue(ctx, config.UserRoleKey, user.Role)
			ctx = context.WithValue(ctx, config.UserMailID, user.Email)
		} else {
			return ctx
		}

		if err != nil {
			color.Red("\nLogin failed: %v. Retry", err)
			continue
		}
		break
	}
	config.CurrentUser = &user

	fmt.Printf(config.WelcomeBack, user.Username, user.Role)
	return ctx
}

func SignupFlow(ctx context.Context) context.Context {
	userRepo := userrepository.NewUserRepository()
	authService := authorisation.NewAuthService(userRepo)

	var username, email, phoneNumber, password string
	fmt.Print("Enter Username: ")
	fmt.Scan(&username)
	for {
		fmt.Print("Enter Email: ")
		fmt.Scan(&email)

		if !(isValidEmail(email)) {
			color.Red("Please enter a valid email id")
			continue
		}
		break
	}

	if authService.UserExists(email) {
		color.Red("Email is already registered. Try logging in.")
		return ctx
	}

	for {
		fmt.Print("Enter Phone number: ")
		fmt.Scan(&phoneNumber)
		if !isValidPhoneNumber(phoneNumber) {
			color.Red("Phone numbers should start with +91 followed by 10 numbers")
			continue
		}
		break
	}
	for {
		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			color.Red("\nError reading password:", err)
			return ctx
		}
		password = string(bytePassword)
		if !isValidPassword(password) {
			color.Red("\nPassword should be 12 characters long with a one special character atleast")
			continue
		}
		break
	}
	hashedPassword, err := authService.HashPassword(password)
	if err != nil {
		color.Red("\nError hashing password:", err)
		return ctx
	}

	newUser := models.User{
		UserID:      uuid.New().String(),
		Username:    username,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    hashedPassword,
		Role:        models.Customer,
		IsBlocked:   false,
	}

	if err := authService.UserRepo.AddUser(newUser); err != nil {
		color.Red("\nFailed to add user:", err)
		return ctx
	}

	ctx = context.WithValue(ctx, config.UserIDKey, newUser.UserID)
	ctx = context.WithValue(ctx, config.UserRoleKey, newUser.Role)
	ctx = context.WithValue(ctx, config.UserMailID, newUser.Email)

	color.Green("\nRegistration successful! You are registered as a %s.\n", newUser.Role)

	return ctx
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	// parts := strings.Split(email, "@")
	// if len(parts) != 2 {
	// 	return false
	// }
	// domain := parts[1]
	// mxRecords, err := net.LookupMX(domain)
	// if err != nil || len(mxRecords) == 0 {
	// 	return false
	// }
	return true
}
func isValidPassword(password string) bool {
	if len(password) < 12 {
		return false
	}
	specialCharRegex := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':"\\|,.<>\/?]`)
	return specialCharRegex.MatchString(password)
}
func isValidPhoneNumber(number string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(number)
}
