package main

import (
	"context"
	"eventro2/config"
	"eventro2/controllers"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	userrepository "eventro2/repository/user_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/bookingservice"
	"eventro2/services/eventservice"
	privilegeservice "eventro2/services/priviligeservice"
	"eventro2/services/searchevents"
	"eventro2/services/showservice"
	"eventro2/services/userservice"
	"eventro2/services/venueservice"
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
	showRepo := showrepository.NewShowRepository()
	eventRepo := eventsrepository.NewEventRepository()
	bookingRepo := bookingrepository.NewBoookingStore()
	userRepo := userrepository.NewUserRepository()
	venueRepo := venuerepository.NewVenueRepository()
	searchService := searchevents.NewSearchService(*eventRepo)
	bookingService := bookingservice.NewBookingService(*bookingRepo, *showRepo)
	showService := showservice.NewShowService(*showRepo, *venueRepo)
	eventService := eventservice.NewEventService(*eventRepo)
	priviligeService := privilegeservice.NewPrivilegeService(*userRepo)
	userService := userservice.NewUserService(*userRepo)
	venueService := venueservice.NewVenueService(*venueRepo)

	role := config.GetUserRole(ctx)
	switch role {
	case models.Admin:
		adminController := controllers.NewAdminController(*priviligeService, *eventService, *userService, *showService, *searchService)
		adminController.ShowAdminDashboard(ctx)
	case models.Host:
		hostController := controllers.NewHostController(*showService, *venueService)
		hostController.ShowHostDashboard(ctx)
	case models.Customer:
		customerController := controllers.NewCustomerController(*searchService, *bookingService)
		customerController.ShowCustomerDashboard(ctx)
	default:
		fmt.Println(config.AccessMessage)
	}
}
