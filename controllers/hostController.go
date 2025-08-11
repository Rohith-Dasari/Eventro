package controllers

import (
	"context"
	"eventro2/config"
	"eventro2/services/showservice"
	"eventro2/services/venueservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type HostController struct {
	showservice.ShowService
	venueservice.VenueService
}

func NewHostController(ss showservice.ShowService, v venueservice.VenueService) *HostController {
	return &HostController{ss, v}
}
func (hc *HostController) ShowHostDashboard(ctx context.Context) {

	for {
		fmt.Println(config.HostDashboard)
		fmt.Println("1. Venue Management")
		fmt.Println("2. Show Management")
		fmt.Println("3. Logout")

		fmt.Println(config.ChoiceMessage)

		choice, _ := utils.TakeUserInput()

		switch choice {
		case 1:
			hc.venueManagement(ctx)
		case 2:
			hc.showManagement(ctx)
		case 3:
			fmt.Println(config.LogoutMessage)
			os.Exit(0)
		}
	}
}

func (hc *HostController) showManagement(ctx context.Context) {
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
			hc.ShowService.BrowseHostShows(ctx)
		case 2:
			hc.ShowService.CreateShow(ctx)
		case 3:
			hc.ShowService.RemoveHostShow(ctx)
		case 4:
			return
		}
	}
}
func (hc *HostController) venueManagement(ctx context.Context) {
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
			hc.VenueService.BrowseHostVenues(ctx)
		case 2:
			hc.VenueService.AddVenue(ctx)
		case 3:
			hc.VenueService.RemoveVenue(ctx)
		case 4:
			return
		}
	}
}
