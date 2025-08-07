package priviligeservice

import (
	"context"
	"eventro/models"
	"eventro/storage"
	"fmt"
)

func EscalatePrivilige(ctx context.Context) {
	fmt.Println("Enter user mail id whose privilige you want to escalate ")
	var usermail string
	fmt.Scanf("%s", &usermail)
	users := storage.LoadUsers()
	var requiredUser *models.User
	for _, user := range users {
		if user.Email == usermail {
			requiredUser = &user
		}
	}
	fmt.Println("are you sure you want to escalate privilige?: y/n")
	var choice string
	fmt.Scanf("%s", &choice)
	if choice == "y" {
		requiredUser.Role = models.Host
		//remove request
	}
}

func SendEscalationRequest(ctx context.Context) {
	//load requests file
	//create a request object

}

func ViewEsclateRequest(ctx context.Context) {
	//see requests made by user
	//include these 2 if we are going with user sending the req

}
