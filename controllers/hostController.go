package controllers

import (
	"context"
	"eventro2/config"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/showservice"
	"eventro2/services/venueservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ShowHostDashboard(ctx context.Context) {

	for {
		fmt.Println(config.HostDashboard)
		fmt.Println("1. Venue Management")
		fmt.Println("2. Show Management")
		fmt.Println("3. Logout")

		fmt.Println(config.ChoiceMessage)

		choice, _ := utils.TakeUserInput()

		switch choice {
		case 1:
			venueManagement(ctx)
		case 2:
			showManagement(ctx)
		case 3:
			fmt.Println("Logging out..")
			os.Exit(0)
		}
	}
}
func showManagement(ctx context.Context) {
	showRepo := showrepository.NewShowRepository()
	venueRepo := venuerepository.NewVenueRepository()
	showService := showservice.NewShowService(*showRepo, *venueRepo)
	for {
		color.Blue("Manage your shows")
		fmt.Println("1. See Shows")
		fmt.Println("2. Create a show")
		fmt.Println("3. Remove Shows")
		fmt.Println("4. Back")
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			showService.BrowseHostShows(ctx)
		case 2:
			showService.CreateShow(ctx)
		case 3:
			showService.RemoveHostShow(ctx)
		case 4:
			return
		}
	}
}
func venueManagement(ctx context.Context) {
	venueRepo := venuerepository.NewVenueRepository()
	venueService := venueservice.NewVenueService(*venueRepo)
	for {
		color.Blue("Mange your venues")
		fmt.Println("1. See Venues")
		fmt.Println("2. Add a Venue")
		fmt.Println("3. Remove Venue")
		fmt.Println("4. Back")
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			venueService.BrowseHostVenues(ctx)
		case 2:
			venueService.AddVenue(ctx)
		case 3:
			venueService.RemoveVenue(ctx)
		case 4:
			return
		}
	}
}
