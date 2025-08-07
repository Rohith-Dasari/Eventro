package userservice

import (
	"context"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	utils "eventro2/utils/userinput"
	"fmt"
)

type UserService struct {
	UserRepo userrepository.UserRepository
}

func NewUserService(userRepo userrepository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (u *UserService) ModerateUser(ctx context.Context) {
	fmt.Println("Enter email ID of user to be moderated:")
	var userMailID string
	fmt.Scanf("%s", &userMailID)

	users := u.UserRepo.Users
	var requiredUser *models.User
	var found bool

	for i := range users {
		if users[i].Email == userMailID {
			requiredUser = &users[i]
			u.PrintUser(users[i])
			found = true
			break
		}
	}

	if !found {
		fmt.Println("User not found, please enter correct ID.")
		return
	}

	if requiredUser.IsBlocked {
		fmt.Println("Are you sure you want to UNBLOCK the user?")
		fmt.Println("1. Yes")
		fmt.Println("2. No")

		choice, err := utils.TakeUserInput()
		if err != nil || (choice != 1 && choice != 2) {
			fmt.Println("Invalid choice. Aborting.")
			return
		}

		if choice == 1 {
			requiredUser.IsBlocked = false
			fmt.Println("User unblocked.")
		} else {
			fmt.Println("Action canceled.")
		}
	} else {
		fmt.Println("Are you sure you want to BLOCK the user?")
		fmt.Println("1. Yes")
		fmt.Println("2. No")

		var choice int
		for {
			var err error
			choice, err = utils.TakeUserInput()
			if err != nil || (choice != 1 && choice != 2) {
				fmt.Println("Invalid choice. Please enter 1 (Yes) or 2 (No).")
				continue
			}
			break
		}

		if choice == 1 {
			requiredUser.IsBlocked = true
			fmt.Println("User blocked.")
		} else {
			fmt.Println("Action canceled.")
		}
	}
	u.UserRepo.SaveUsers(users)
}

func (u *UserService) ViewBlockedUsers(ctx context.Context) {
	users := u.UserRepo.Users
	for _, user := range users {
		if user.IsBlocked {
			u.PrintUser(user)
		}
	}

}

func (u *UserService) PrintUser(user models.User) {
	fmt.Println("Username: ", user.Username)
	fmt.Println("Email: ", user.Email)
	fmt.Println("Phone Number: ", user.PhoneNumber)
}
