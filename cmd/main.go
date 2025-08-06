package main

import (
	"context"
	"eventro/config"
	"eventro/controllers"
	"eventro/models"
	"fmt"
)

func main() {
	fmt.Println(config.WelcomeMessage)
	ctx := context.Background()
	ctx = startApp(ctx)
	if config.CurrentUser == nil {
		fmt.Println("Login or Signup failed. Exiting.")
		return
	}
	launchDashboard(ctx)
}

func startApp(ctx context.Context) context.Context {
	for {
		fmt.Println("1. Login\n2. Signup")
		fmt.Print(config.ChoiceMessage)
		var choice int
		fmt.Scanf("%d", &choice)
		switch choice {
		case 1:
			return controllers.LoginFlow(ctx)
		case 2:
			return controllers.SignupFlow(ctx)
		default:
			fmt.Println("Please enter 1 or 2")
		}
	}
}

func launchDashboard(ctx context.Context) {
	role := config.GetUserRole(ctx)
	switch role {
	case models.Admin:
		controllers.ShowAdminDashboard(ctx)
	case models.Host:
		controllers.ShowHostDashboard(ctx)
	case models.Customer:
		controllers.ShowCustomerDashboard(ctx)
	default:
		fmt.Println("Access Denied")
	}
}
