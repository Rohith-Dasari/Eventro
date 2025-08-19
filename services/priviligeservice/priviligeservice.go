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
	UserRepo userrepository.UserRepository
}

func NewPrivilegeService(repo userrepository.UserRepository) PrivilegeService {
	return PrivilegeService{
		UserRepo: repo,
	}
}

func (p *PrivilegeService) EscalatePrivilege(ctx context.Context) {
	for {
		color.Blue("Event Moderation")
		fmt.Print("Enter user email ID to change privilege: ")
		email := utils.ReadLine()

		user, err := p.UserRepo.GetByEmail(email)
		if err != nil || user == nil {
			color.Red("User not found.")
			fmt.Println("1. Retry with another User ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		// If already Host
		if user.Role == models.Host {
			color.Red("User is currently a Host")
			fmt.Println("1. Change to Customer \n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				fmt.Println("Are you sure you want to downgrade privilege to Customer?\n1. Yes \n2. No")
				fmt.Println(config.ChoiceMessage)
				if utils.ReadLine() == "1" {
					user.Role = models.Customer
					if err := p.UserRepo.Update(user); err != nil {
						color.Red("Failed to update user: %v", err)
						return
					}
					color.Green("User privilege changed to Customer successfully.")
					return
				}
				color.Red("Privilege change canceled.")
				return
			}
			return
		}

		fmt.Println("Entered User's role:", user.Role)
		fmt.Println("Are you sure you want to escalate privilege to HOST?\n1. Yes \n2. No:")
		if utils.ReadLine() == "1" {
			user.Role = models.Host
			if err := p.UserRepo.Update(user); err != nil {
				color.Red("Failed to update user: %v", err)
				return
			}
			color.Green("User privilege escalated to Host successfully.")
			return
		}

		color.Red("Privilege escalation canceled.")
		return
	}
}
