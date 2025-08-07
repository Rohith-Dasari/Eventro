package controllers

import (
	"bufio"
	"context"
	"eventro2/config"
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
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"strings"
)

func ShowAdminDashboard(ctx context.Context) {
	fmt.Println("\n=== Admin Dashboard ===")
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
		fmt.Println("1. User Moderation 2. Show Moderation 3. Event Moderation 4. Add a new host 5. Add a new Event 6. Book behalf of user")

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
			event := eventService.CreateNewEvent()
			events := eventService.EventRepo.Events
			events = append(events, event)
			eventService.EventRepo.SaveEvents(events)
		case 6:
			bookBehalfofUser(ctx)
		}
	}
}

func userModeration(ctx context.Context) {
	UserRepo := userrepository.NewUserRepository()
	userService := userservice.NewUserService(*UserRepo)
	fmt.Println("1. block/unlock user 2. view blocked user ")
	fmt.Println(config.ChoiceMessage)
	choice, _ := utils.TakeUserInput()
	switch choice {
	case 1:
		userService.ModerateUser(ctx)
	case 2:
		userService.ViewBlockedUsers(ctx)
	}
}
func showModeration(ctx context.Context) {
	showRepo := showrepository.NewShowRepository()
	venueRepo := venuerepository.NewVenueRepository()
	showService := showservice.NewShowService(*showRepo, *venueRepo)
	fmt.Println("1. block/unblock show 2. view blocked show")
	fmt.Println(config.ChoiceMessage)
	choice, _ := utils.TakeUserInput()
	switch choice {
	case 1:
		showService.ModerateShow(ctx)
	case 2:
		showService.ViewBlockedShows(ctx)
	}
}
func eventModeration(ctx context.Context) {
	eventRepo := eventsrepository.NewEventRepository()
	eventService := eventservice.NewEventService(*eventRepo)
	fmt.Println("1. block/unblock event 2. view blocked event")
	fmt.Println(config.ChoiceMessage)
	choice, _ := utils.TakeUserInput()
	switch choice {
	case 1:
		eventService.ModerateEvents(ctx)
	case 2:
		eventService.ViewBlockedEvents(ctx)
	}
}
func bookBehalfofUser(ctx context.Context) {
	showRepo := showrepository.NewShowRepository()
	venueRepo := venuerepository.NewVenueRepository()
	bookRepo := bookingrepository.NewBoookingStore()
	bookingService := bookingservice.NewBookingService(*bookRepo, *showRepo)
	showService := showservice.NewShowService(*showRepo, *venueRepo)
	eventRepo := eventsrepository.NewEventRepository()
	searchEvents := searchevents.NewSearchService(*eventRepo)
	reader := bufio.NewReader(os.Stdin)
	searchEvents.Search(ctx)
	var eventID string
	fmt.Println("Enter event id:")
	eventID, _ = reader.ReadString('\n')
	eventID = strings.TrimSpace(eventID)
	showService.BrowseShowsByEvent(ctx, eventID)
	var showID string
	fmt.Println("enter the showID you want to pick")
	showID, _ = reader.ReadString('\n')
	showID = strings.TrimSpace(showID)
	showService.DisplayShow(ctx, showID)
	fmt.Println("enter the userID of user you want to pick")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)
	bookingService.MakeBooking(ctx, userID, showID)
}
