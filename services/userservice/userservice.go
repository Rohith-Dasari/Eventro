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

func NewUserService(userRepo userrepository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (u *UserService) ModerateUser(ctx context.Context) {
	for {
		fmt.Println("Enter email ID of user to be moderated:")
		userMailID := utils.ReadLine()

		var requiredUser *models.User
		var found bool
		var requiredIndex int

		for i := range u.UserRepo.Users {
			if u.UserRepo.Users[i].Email == userMailID {
				requiredUser = &u.UserRepo.Users[i]
				u.PrintUser(u.UserRepo.Users[i])
				found = true
				requiredIndex = i
				break
			}
		}

		if !found {
			color.Red("User not found, please enter correct ID.")
			fmt.Println("1. Retry with another Email ID\n2. Back ")
			fmt.Println(config.ChoiceMessage)
			choice, err := utils.TakeUserInput()
			if err != nil {
				fmt.Println("please enter an integer")
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

		if requiredUser.IsBlocked {
			fmt.Println("The user is currently BLOCKED")
			fmt.Println("Are you sure you want to UNBLOCK the user?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				requiredUser.IsBlocked = !requiredUser.IsBlocked
				u.UserRepo.Users[requiredIndex].IsBlocked = false
				u.UserRepo.SaveUsers(u.UserRepo.Users)
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
				requiredUser.IsBlocked = !requiredUser.IsBlocked
				u.UserRepo.Users[requiredIndex].IsBlocked = true
				u.UserRepo.SaveUsers(u.UserRepo.Users)
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
	users := u.UserRepo.Users
	found := false
	color.Blue("Blocked Users:")
	for _, user := range users {
		if user.IsBlocked {
			u.PrintUser(user)
			found = true
		}
	}
	if !found {
		color.Red("No Blocked Users ")
	}

}

func (u *UserService) PrintUser(user models.User) {
	fmt.Println(config.Dash)
	fmt.Println("Username: ", user.Username)
	fmt.Println("Email: ", user.Email)
	fmt.Println("Phone Number: ", user.PhoneNumber)
	fmt.Println(config.Dash)
}
