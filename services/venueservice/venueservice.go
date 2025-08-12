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
	VenueRepo venuerepository.VenueStorageI
}

func NewVenueService(repo venuerepository.VenueStorageI) *VenueService {
	return &VenueService{VenueRepo: repo}
}

func (vs *VenueService) AddVenue(ctx context.Context) {

	var name, city, state, seatLayoutInput string
	for {
		fmt.Print("Enter Venue Name: ")
		name = utils.ReadLine()
		if name == "" {
			color.Red("Venue name cannot be empty.")
			continue
		}
		break
	}

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
	fmt.Print("Enter State: ")
	state = utils.ReadLine()

	fmt.Print("Is seat layout required? (y/n): ")
	seatLayoutInput = utils.ReadLine()

	seatLayoutInput = strings.ToLower(strings.TrimSpace(seatLayoutInput))

	isSeatLayoutRequired := seatLayoutInput == "y"

	venue := models.Venue{
		ID:                   uuid.New().String(),
		Name:                 name,
		HostID:               config.GetUserID(ctx),
		City:                 city,
		State:                state,
		IsSeatLayoutRequired: isSeatLayoutRequired,
	}

	venues, _ := vs.VenueRepo.GetVenues()
	venues = append(venues, venue)
	vs.VenueRepo.SaveVenues(venues)
	fmt.Println("Venue created successfully:")
	vs.PrintVenue(venue)
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
	venues, _ := vs.VenueRepo.GetVenues()
	color.Blue("Your Venues")
	for _, venue := range venues {
		if venue.HostID == config.GetUserID(ctx) {
			vs.PrintVenue(venue)
		}
	}
}

func (vs *VenueService) RemoveVenue(ctx context.Context) {
	for {
		fmt.Print("Enter Venue ID to remove: ")
		venueID := utils.ReadLine()

		venues, _ := vs.VenueRepo.GetVenues()

		var requiredVenue models.Venue
		var requiredIndex int
		var found bool
		for i, venue := range venues {
			if venue.ID == venueID {
				requiredVenue = venue
				requiredIndex = i
				found = true
			}
		}
		if !found || requiredVenue.HostID != config.GetUserID(ctx) {
			color.Red("You can only remove your venues")
			fmt.Println("1. Retry another Venue ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			case 2:
				return
			default:
				fmt.Println(config.InvalidMSG)
			}
		} else {

			vs.PrintVenue(requiredVenue)
		loop:
			for {
				fmt.Println("Are you sure you want to delete this venue?")
				fmt.Println("1. Yes")
				fmt.Println("2. No, try another ID")
				fmt.Println("3. Back")
				fmt.Println(config.ChoiceMessage)
				choice, _ := utils.TakeUserInput()
				switch choice {
				case 1:
					venues = append(venues[:requiredIndex], venues[requiredIndex+1:]...)
					vs.VenueRepo.SaveVenues(venues)
					fmt.Println("Venue removed successfully.")
					return
				case 2:
					color.Red("Cancelled.")

					break loop
				case 3:
					return
				default:
					fmt.Println(config.InvalidMSG)
				}
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
