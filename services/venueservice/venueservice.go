package venueservice

import (
	"context"
	"eventro2/config"
	"eventro2/models"
	venuerepository "eventro2/repository/venue_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type VenueService struct {
	VenueRepo venuerepository.VenueRepository
}

func NewVenueService(repo venuerepository.VenueRepository) VenueService {
	return VenueService{VenueRepo: repo}
}

func (vs *VenueService) AddVenue(ctx context.Context) {
	var name, city, state, seatLayoutInput string

	// Venue Name
	for {
		fmt.Print("Enter Venue Name: ")
		name = utils.ReadLine()
		if name == "" {
			color.Red("Venue name cannot be empty.")
			continue
		}
		break
	}

	// City
	for {
		fmt.Print("Enter City: ")
		city = utils.ReadLine()
		if city == "" {
			color.Red("City cannot be empty.")
			continue
		}
		if !isAlpha(city) {
			color.Red("City can only contain letters and spaces.")
			continue
		}
		break
	}

	// State
	fmt.Print("Enter State: ")
	state = utils.ReadLine()

	// Seat layout requirement
	fmt.Print("Is seat layout required? (y/n): ")
	seatLayoutInput = utils.ReadLine()
	seatLayoutInput = strings.ToLower(strings.TrimSpace(seatLayoutInput))
	isSeatLayoutRequired := seatLayoutInput == "y"

	// Build Venue
	venue := &models.Venue{
		ID:                   uuid.New().String(),
		Name:                 name,
		HostID:               config.GetUserID(ctx),
		City:                 city,
		State:                state,
		IsSeatLayoutRequired: isSeatLayoutRequired,
	}

	// Save directly via repo
	if err := vs.VenueRepo.Create(venue); err != nil {
		color.Red("Failed to create venue: %v", err)
		return
	}

	color.Green("Venue created successfully:")
	vs.PrintVenue(*venue)
}

func isAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
func (vs *VenueService) BrowseHostVenues(ctx context.Context) {
	hostID := config.GetUserID(ctx)

	venues, err := vs.VenueRepo.ListByHost(hostID)
	if err != nil {
		color.Red("Failed to fetch venues: %v", err)
		return
	}

	if len(venues) == 0 {
		color.Yellow("You have no venues.")
		return
	}

	color.Blue("Your Venues:")
	for _, venue := range venues {
		vs.PrintVenue(venue)
	}
}
func (vs *VenueService) RemoveVenue(ctx context.Context) {
	for {
		fmt.Print("Enter Venue ID to remove: ")
		venueID := utils.ReadLine()

		venue, err := vs.VenueRepo.GetByID(venueID)
		if err != nil || venue == nil {
			color.Red("Venue not found.")
			fmt.Println("1. Retry another Venue ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		// Check ownership
		if venue.HostID != config.GetUserID(ctx) {
			color.Red("You can only remove your venues.")
			fmt.Println("1. Retry another Venue ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		vs.PrintVenue(*venue)

		for {
			fmt.Println("Are you sure you want to delete this venue?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			fmt.Println(config.ChoiceMessage)

			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				if err := vs.VenueRepo.Delete(venueID); err != nil {
					color.Red("Failed to remove venue: %v", err)
				} else {
					color.Green("Venue removed successfully.")
				}
				return
			case 2:
				color.Red("Cancelled.")
				break
			case 3:
				return
			default:
				fmt.Println(config.InvalidMSG)
			}
		}
	}
}

func (vs *VenueService) PrintVenue(venue models.Venue) {
	fmt.Println(config.Dash)
	fmt.Println("ID:                    ", venue.ID)
	fmt.Println("Name:                  ", venue.Name)
	fmt.Println("City:                  ", venue.City)
	fmt.Println("State:                 ", venue.State)
	fmt.Println("Seat Layout Required?: ", venue.IsSeatLayoutRequired)
	fmt.Println(config.Dash)
}
