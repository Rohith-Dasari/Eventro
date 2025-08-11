package controllers

import (
	"context"
	"eventro2/config"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	userrepository "eventro2/repository/user_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/eventservice"
	privilegeservice "eventro2/services/priviligeservice"
	"eventro2/services/searchevents"
	"eventro2/services/showservice"
	"eventro2/services/userservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ShowAdminDashboard(ctx context.Context) {

	// bookRepo := bookingrepository.NewBoookingStore()
	// showRepo := showrepository.NewShowRepository()
	UserRepo := userrepository.NewUserRepository()
	eventRepo := eventsrepository.NewEventRepository()
	//bookService := bookingservice.NewBookingService(*bookRepo, *showRepo)
	privilegeService := privilegeservice.NewPrivilegeService(*UserRepo)
	eventService := eventservice.NewEventService(*eventRepo)

	//view users, block users, unblock users, assign new host--done
	//view show, block show/event, unblock show/event,-done
	//book a show behalf of user
	//create privilige service
	//create user service first
	//create event service
	//create show service
	//create search service
	//create bboking
	for {
		fmt.Println(config.AdminDashboard)
		fmt.Println("1. User Moderation\n2. Show Moderation\n3. Event Moderation\n4. Privilige Moderation\n5. Add a new Event\n6. Book behalf of user\n7. Logout")

		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()

		switch choice {
		case 1:
			userModeration(ctx)
		case 2:
			showModeration(ctx)
		case 3:
			eventModeration(ctx)
		case 4:
			privilegeService.EscalatePrivilege(ctx)
		case 5:
			eventService.CreateNewEvent()

		case 6:
			bookBehalfofUser(ctx)
		case 7:
			fmt.Println("Logging out..")
			os.Exit(0)
		}
	}
}

func userModeration(ctx context.Context) {
	UserRepo := userrepository.NewUserRepository()
	userService := userservice.NewUserService(*UserRepo)

	for {
		color.Blue("User Moderation")
		fmt.Println("1. Block/Unblock User\n2. View Blocked Users\n3. Back")
		fmt.Println(config.ChoiceMessage)
		choice, err := utils.TakeUserInput()
		if err != nil {
			color.Red(err.Error())
			continue
		}
		switch choice {
		case 1:
			userService.ModerateUser(ctx)
		case 2:
			userService.ViewBlockedUsers(ctx)
		case 3:
			return
		default:
			fmt.Println(config.InvalidMSG)
		}
	}
}
func showModeration(ctx context.Context) {
	showRepo := showrepository.NewShowRepository()
	venueRepo := venuerepository.NewVenueRepository()
	showService := showservice.NewShowService(*showRepo, *venueRepo)
	for {
		color.Blue("Show Moderation")
		fmt.Println(config.ShowModerationMSG)
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			showService.ModerateShow(ctx)
		case 2:
			showService.ViewBlockedShows(ctx)
		case 3:
			return
		}
	}
}
func eventModeration(ctx context.Context) {
	eventRepo := eventsrepository.NewEventRepository()
	eventService := eventservice.NewEventService(*eventRepo)
	for {
		color.Blue("Event Moderation")
		fmt.Println(config.EventModeration)
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			eventService.ModerateEvents(ctx)
		case 2:
			eventService.ViewBlockedEvents(ctx)
		case 3:
			return
		default:
			fmt.Println(config.InvalidMSG)
		}
	}
}
func bookBehalfofUser(ctx context.Context) {
	eventRepo := eventsrepository.NewEventRepository()
	searchEvents := searchevents.NewSearchService(*eventRepo)
	searchEvents.Search(ctx)
}
