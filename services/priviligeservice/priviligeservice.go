package privilegeservice

import (
	"context"
	"eventro2/config"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	utils "eventro2/utils/userinput"
	"fmt"

	"github.com/fatih/color"
)

type PrivilegeService struct {
	UserRepo userrepository.UserStorageI
}

func NewPrivilegeService(repo userrepository.UserStorageI) *PrivilegeService {
	return &PrivilegeService{
		UserRepo: repo,
	}
}

func (p *PrivilegeService) EscalatePrivilege(ctx context.Context) {
	for {
		color.Blue("Event Moderation")
		fmt.Print("Enter user email ID to change privilege: ")
		email := utils.ReadLine()

		users, _ := p.UserRepo.GetUsers()
		var targetUser *models.User
		var requiredIndex int

		for i := range users {
			if users[i].Email == email {
				targetUser = &users[i]
				requiredIndex = i
				break
			}
		}

		if targetUser == nil {
			color.Red("User not found.")
			fmt.Println("1. Retry with another User ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			default:
				return
			}
		}
		if targetUser.Role == models.Host {
			color.Red("User is a host")
			fmt.Println("1. Change to Customer \n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:

				fmt.Println("Are you sure you want to change privilege to customer?\n1. Yes \n2. No:")
				fmt.Println(config.ChoiceMessage)
				choice := utils.ReadLine()

				if choice == "1" {
					targetUser.Role = models.Customer
					users[requiredIndex].Role = models.Customer
					p.UserRepo.SaveUsers(users)
					color.Green("User privilege changed to Customer successfully.")
					return
				} else {
					color.Red("Privilege escalation canceled.")
					return
				}
			default:
				return
			}
		}
		fmt.Println("Entered User ID is of role: ", users[requiredIndex].Role)

		fmt.Println("Are you sure you want to escalate privilege to HOST?\n1. Yes \n2. No:")
		fmt.Println(config.ChoiceMessage)
		choice := utils.ReadLine()

		if choice == "1" {
			targetUser.Role = models.Host
			users[requiredIndex].Role = models.Admin
			color.Green("User privilege escalated successfully.")
			p.UserRepo.SaveUsers(users)
			return
		} else {
			color.Red("Privilege escalation canceled.")
			return
		}
	}
}
