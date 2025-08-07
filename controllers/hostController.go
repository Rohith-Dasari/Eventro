package controllers

import (
	"context"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/showservice"
	"eventro2/services/venueservice"
	utils "eventro2/utils/userinput"
	"fmt"
)

func ShowHostDashboard(ctx context.Context) {
	fmt.Println("\n=== Host Dashboard ===")
	fmt.Println("1. Create a show")
	fmt.Println("2. See Shows")
	fmt.Println("3. Remove Shows")
	fmt.Println("4. Add a Venue")
	fmt.Println("5. See Venues")
	fmt.Println("6. Remove Venue")

	choice, _ := utils.TakeUserInput()
	showRepo := showrepository.NewShowRepository()
	venueRepo := venuerepository.NewVenueRepository()
	showService := showservice.NewShowService(*showRepo, *venueRepo)
	venueService := venueservice.NewVenueService(*venueRepo)

	switch choice {
	case 1:
		showService.CreateShow(ctx)
	case 2:
		showService.BrowseHostShows(ctx)
	case 3:
		showService.RemoveHostShow(ctx)
	case 4:
		venueService.AddVenue(ctx)
	case 5:
		venueService.BrowseHostVenues(ctx)
	case 6:
		venueService.RemoveVenue(ctx)
	default:
		fmt.Println("Invalid choice.")
	}
}
