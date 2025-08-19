package userservice

import (
	"context"
	"eventro2/config"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	utils "eventro2/utils/userinput"
	"fmt"

	"github.com/fatih/color"
)

type UserService struct {
	UserRepo userrepository.UserRepository
}

func NewUserService(userRepo userrepository.UserRepository) UserService {
	return UserService{UserRepo: userRepo}
}

func (u *UserService) ModerateUser(ctx context.Context) {
	for {
		fmt.Println("Enter email ID of user to be moderated:")
		userMailID := utils.ReadLine()

		requiredUser, err := u.UserRepo.GetByEmail(userMailID)
		if err != nil || requiredUser == nil {
			color.Red("User not found, please enter correct ID.")
			fmt.Println("1. Retry with another Email ID\n2. Back ")
			fmt.Println(config.ChoiceMessage)
			choice, err := utils.TakeUserInput()
			if err != nil {
				color.Red("please enter an integer")
				continue
			}

			switch choice {
			case 1:
				continue
			case 2:
				return
			default:
				return
			}
		}

		u.PrintUser(*requiredUser)

		if requiredUser.IsBlocked {
			fmt.Println("The user is currently BLOCKED")
			fmt.Println("Are you sure you want to UNBLOCK the user?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				requiredUser.IsBlocked = false
				if err := u.UserRepo.Update(requiredUser); err != nil {
					color.Red("Failed to update user: %v", err)
					return
				}
				color.Green("User is successfully unblocked")
				return
			case 2:
				color.Red("Canceled")
				continue
			case 3:
				return
			}
		} else {
			fmt.Println("The User is currently UNBLOCKED")
			fmt.Println("Are you sure you want to BLOCK the User?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				requiredUser.IsBlocked = true
				if err := u.UserRepo.Update(requiredUser); err != nil {
					color.Red("Failed to update user: %v", err)
					return
				}
				color.Green("User is successfully blocked")
				return
			case 2:
				color.Red("Canceled")
				continue
			case 3:
				return
			}
		}
	}
}

func (u *UserService) ViewBlockedUsers(ctx context.Context) {
	users, err := u.UserRepo.GetBlockedUsers()
	if err != nil {
		color.Red("Error fetching blocked users: %v", err)
		return
	}

	if len(users) == 0 {
		color.Red("No Blocked Users")
		return
	}

	color.Blue("Blocked Users:")
	for _, user := range users {
		u.PrintUser(user)
	}
}

func (u *UserService) PrintUser(user models.User) {
	fmt.Println(config.Dash)
	fmt.Printf("%-14s : %s\n", "Username", user.Username)
	fmt.Printf("%-14s : %s\n", "Email", user.Email)
	fmt.Printf("%-14s : %s\n", "Phone Number", user.PhoneNumber)
	fmt.Println(config.Dash)
}
