package controllers

import (
	"bufio"
	"context"
	"eventro2/config"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	"eventro2/services/bookingservice"
	"eventro2/services/searchevents"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ShowCustomerDashboard(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	showRepo := showrepository.NewShowRepository()
	//venueRepo := venuerepository.NewVenueRepository()
	eventRepo := eventsrepository.NewEventRepository()
	bookingRepo := bookingrepository.NewBoookingStore()

	//showService := showservice.NewShowService(*showRepo, *venueRepo)
	searchService := searchevents.NewSearchService(*eventRepo)
	bookingService := bookingservice.NewBookingService(*bookingRepo, *showRepo)

	for {
		fmt.Println(config.CustomerDashboardMessage)
		fmt.Println(config.CustomerOptions)
		var choice int
		fmt.Println(config.ChoiceMessage)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, _ = strconv.Atoi(input)

		switch choice {
		case 1:
			searchService.Search(ctx)

		case 2:
			bookingService.ViewBookingHistory(ctx)

		case 3:
			fmt.Println("Logging out...")
			return

		default:
			fmt.Println(config.DefaultChoiceMessage)
		}
	}
}
