package privilegeservice

import (
	"context"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"fmt"
)

type PrivilegeService struct {
	UserRepo userrepository.UserRepository
}

func NewPrivilegeService(repo userrepository.UserRepository) *PrivilegeService {
	return &PrivilegeService{
		UserRepo: repo,
	}
}

func (p *PrivilegeService) EscalatePrivilege(ctx context.Context) {
	fmt.Print("Enter user email ID to escalate privilege: ")
	var email string
	fmt.Scanf("%s", &email)

	users := p.UserRepo.Users
	var targetUser *models.User

	for i := range users {
		if users[i].Email == email {
			targetUser = &users[i]
			break
		}
	}

	if targetUser == nil {
		fmt.Println("User not found.")
		return
	}

	fmt.Println("Are you sure you want to escalate privilege to HOST? (y/n): ")
	var choice string
	fmt.Scanf("%s", &choice)

	if choice == "y" {
		targetUser.Role = models.Host
		fmt.Println("User privilege escalated successfully.")
		p.UserRepo.SaveUsers(users)
	} else {
		fmt.Println("Privilege escalation canceled.")
	}
}
