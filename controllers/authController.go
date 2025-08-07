package controllers

import (
	"bufio"
	"context"
	"eventro/config"
	"eventro/models"
	"eventro/services/authorisation"
	"eventro/storage"
	"fmt"
	"net"
	"net/mail"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/term"
)

func LoginFlow(ctx context.Context) context.Context {

	users := storage.LoadUsers()
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
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return ctx
		}
		password = strings.TrimSpace(input)

		// bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		// if err != nil {
		// 	fmt.Println("\nError reading password:", err)
		// 	return ctx
		// }
		// password = string(bytePassword)

		user, err = authorisation.ValidateLogin(users, email, password)
		ctx = context.WithValue(ctx, config.UserIDKey, user.UserID)
		ctx = context.WithValue(ctx, config.UserRoleKey, user.Role)
		ctx = context.WithValue(ctx, config.UserMailID, user.Email)

		if err != nil {
			fmt.Println("Retry. Login failed:", err)
			continue
		}
		break
	}
	if user.IsBlocked {
		fmt.Println("Account Blocked. Please contact Admin")
		return ctx
	}
	config.CurrentUser = &user

	fmt.Printf("Welcome back %s! You are logged in as: %s\n", user.Username, user.Role)
	return ctx
}

func SignupFlow(ctx context.Context) context.Context {
	users := storage.LoadUsers()

	var username, email, phoneNumber, password string
	fmt.Print("Enter Username: ")
	fmt.Scan(&username)
	fmt.Print("Enter Email: ")
	fmt.Scan(&email)

	if !(isValid(email)) {
		fmt.Println("Please enter a valid email id")
		return ctx
	}

	if storage.UserExists(users, email) {
		fmt.Println("Email is already registered. Try logging in.")
		return ctx
	}

	fmt.Print("Enter Phone number: ")
	fmt.Scan(&phoneNumber)
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("\nError reading password:", err)
		return ctx
	}
	password = string(bytePassword)
	hashedPassword, err := authorisation.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
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

	users = append(users, newUser)
	if err := storage.SaveUsers(users); err != nil {
		fmt.Println("Failed to save user:", err)
		return ctx
	}
	ctx = context.WithValue(ctx, config.UserIDKey, newUser.UserID)
	ctx = context.WithValue(ctx, config.UserRoleKey, newUser.Role)
	ctx = context.WithValue(ctx, config.UserMailID, newUser.Email)

	fmt.Printf("Registration successful! You are registered as a %s.\n", newUser.Role)
	return ctx
}

func isValid(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}
	return true
}
