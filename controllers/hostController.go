package controllers

import (
	"context"
	"eventro/services/showservice"
	"eventro/services/venueservice"
	utils "eventro/utils/userinput"
	"fmt"
)

func ShowHostDashboard(ctx context.Context) {

	fmt.Println("\n=== Host Dashboard ===")
	// future: create event, view bookings, etc.
	//add show, remove show,
	//see host's shows booking info
	//see host's venues
	//add or remove venues
	fmt.Println("1. Create a show 2. See Shows 3. Remove shows 4. Add a venue 5. See venues 6. Remove venue")
	choice, _ := utils.TakeUserInput()
	switch choice {
	case 1:
		showservice.CreateShow(ctx)
	case 2:
		showservice.BrowseHostShows(ctx)
	case 3:
		showservice.RemoveHostShows(ctx)
	case 4:
		venueservice.AddVenue(ctx)
	case 5:
		venueservice.BrowseHostVenues(ctx)
	case 6:
		venueservice.RemoveVenue(ctx)
	}

}
