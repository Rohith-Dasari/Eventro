package userservice

import (
	"context"
	"eventro/models"
	"eventro/storage"
	"fmt"
)

func ModerateUser(ctx context.Context) {
	fmt.Println("enter mailid of user want to block")
	var userMailID string
	fmt.Scanf("%s", &userMailID)
	users := storage.LoadUsers()
	var requiredUser *models.User
	var found bool
	for _, user := range users {
		if user.Email == userMailID {
			requiredUser = &user
			PrintUser(user)
			found = true
		}
	}
	if !found {
		fmt.Println("User not found, please enter correct ID")
	} else {
		if requiredUser.Role != "blocked" {
			fmt.Print("Are you sure you want to unblock the user: y/n")
			requiredUser.isBlocked = true
		} else {
			fmt.Printf("Are you sure you want to block the user: y/n")
			var choice string
			fmt.Scanf("%s", choice)
			if choice == "y" {
				requiredUser.isBlocked = true
			}
		}
	}
}
func ViewBlockedUsers(ctx context.Context) {

}

func PrintUser(user models.User) {
	fmt.Print("Username: ", user.Username)
	fmt.Print("Email: ", user.Email)
	fmt.Print("Phone Number: ", user.PhoneNumber)
}
