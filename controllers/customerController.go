package controllers

import (
	"bufio"
	"context"
	"eventro2/config"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/bookingservice"
	"eventro2/services/searchevents"
	"eventro2/services/showservice"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ShowCustomerDashboard(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	showRepo := showrepository.NewShowRepository()
	venueRepo := venuerepository.NewVenueRepository()
	eventRepo := eventsrepository.NewEventRepository()
	bookingRepo := bookingrepository.NewBoookingStore()

	showService := showservice.NewShowService(*showRepo, *venueRepo)
	searchService := searchevents.NewSearchService(*eventRepo)
	bookingService := bookingservice.NewBookingService(*bookingRepo, *showRepo)

	for {
		fmt.Println(config.CustomerDashboardMessage)
		fmt.Println("1. Search 2. Booking History 3. Logout")
		var choice int
		fmt.Println(config.ChoiceMessage)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, _ = strconv.Atoi(input)

		switch choice {
		case 1:
			searchService.Search(ctx)
			fmt.Println("Enter event ID:")
			eventID, _ := reader.ReadString('\n')
			eventID = strings.TrimSpace(eventID)

			showService.BrowseShowsByEvent(ctx, eventID)

			fmt.Println("Enter show ID:")
			showID, _ := reader.ReadString('\n')
			showID = strings.TrimSpace(showID)

			showService.DisplayShow(ctx, showID)

			bookingService.MakeBooking(ctx, config.GetUserID(ctx), showID)

		case 2:
			bookingService.ViewBookingHistory(ctx)

		case 3:
			fmt.Println("Logging out...")
			return

		default:
			fmt.Println("Invalid choice. Please select from 1 to 3.")
		}
	}
}
