package controllers

import (
	"bufio"
	"context"
	"eventro/config"
	"eventro/services/bookingservice"
	"eventro/services/searchevents"
	"eventro/services/showservice"
	"eventro/storage"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ShowCustomerDashboard(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
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
			searchevents.Search()
			var eventID string
			fmt.Println("Enter event id:")
			reader = bufio.NewReader(os.Stdin)
			eventID, _ = reader.ReadString('\n')
			eventID = strings.TrimSpace(eventID)
			shows := storage.LoadShows()
			showservice.BrowseShowsByEvent(eventID, shows)
			var showID string
			fmt.Println("enter the showID you want to pick")
			showID, _ = reader.ReadString('\n')
			showID = strings.TrimSpace(showID)
			showservice.DisplayShow(showID, shows)
			bookingservice.MakeBooking(config.GetUserID(ctx), showID)
		case 2:
			bookingservice.SeeBookingHistory()
		case 3:
			fmt.Println("Logging out...")
			return
		default:
			fmt.Println("Invalid choice. Please select from 1 to 3.")
		}
	}
}
