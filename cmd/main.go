package main

import (
	"context"
	"eventro2/config"
	"eventro2/controllers"
	"eventro2/models"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
)

func main() {
	fmt.Println(config.WelcomeMessage)
	ctx := context.Background()
	for {
		ctx = startApp(ctx)
		if config.GetUserID(ctx) == "" {
			fmt.Println(config.LoginErrorMessage)
			continue
		}
		break
	}
	launchDashboard(ctx)
}

func startApp(ctx context.Context) context.Context {
	fmt.Println(config.StartAppMessage)
	for {
		fmt.Print(config.ChoiceMessage)

		choice, err := utils.TakeUserInput()
		if err != nil {
			continue
		}
		switch choice {
		case 1:
			return controllers.LoginFlow(ctx)
		case 2:
			return controllers.SignupFlow(ctx)
		case 3:
			fmt.Println(config.LogoutMessage)
			os.Exit(0)
		default:
			fmt.Println(config.DefaultChoiceMessage)
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
		fmt.Println(config.AccessMessage)
	}
}
