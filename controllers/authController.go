package controllers

import (
	"context"
	"eventro/config"
	"eventro/models"
	"eventro/services/authorisation"
	"eventro/storage"
	"fmt"

	"github.com/google/uuid"
)

func LoginFlow(ctx context.Context) context.Context {
	users := storage.LoadUsers()

	var email, password string
	fmt.Print("Enter email: ")
	fmt.Scan(&email)
	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	user, err := authorisation.ValidateLogin(users, email, password)
	ctx = context.WithValue(ctx, config.UserIDKey, user.UserID)
	ctx = context.WithValue(ctx, config.UserRoleKey, user.Role)
	if err != nil || user.IsBlocked {
		fmt.Println("Login failed:", err)
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

	if storage.UserExists(users, email) {
		fmt.Println("Email is already registered. Try logging in.")
		return ctx
	}

	fmt.Print("Enter Phone number: ")
	fmt.Scan(&phoneNumber)
	fmt.Print("Enter Password: ")
	fmt.Scan(&password)

	role := authorisation.AssignRole(email)
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
		Role:        role,
		IsBlocked:   false,
	}

	users = append(users, newUser)
	if err := storage.SaveUsers(users); err != nil {
		fmt.Println("Failed to save user:", err)
		return ctx
	}
	ctx = context.WithValue(ctx, config.UserIDKey, newUser.UserID)
	ctx = context.WithValue(ctx, config.UserRoleKey, newUser.Role)

	fmt.Printf("Registration successful! You are registered as a %s.\n", newUser.Role)
	return ctx
}
