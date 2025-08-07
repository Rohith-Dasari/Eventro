package venueservice

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	venuerepository "eventro2/repository/venue_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

type VenueService struct {
	VenueRepo venuerepository.VenueRepository
}

func NewVenueService(repo venuerepository.VenueRepository) *VenueService {
	return &VenueService{VenueRepo: repo}
}

func (vs *VenueService) AddVenue(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Venue Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Enter City: ")
	city, _ := reader.ReadString('\n')

	fmt.Print("Enter State: ")
	state, _ := reader.ReadString('\n')

	fmt.Print("Is seat layout required? (y/n): ")
	seatLayoutInput, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)
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

	venues := vs.VenueRepo.Venues
	venues = append(venues, venue)
	vs.VenueRepo.SaveVenues(venues)

	fmt.Println("Venue created successfully:")
	vs.PrintVenue(venue)
}

func (vs *VenueService) BrowseHostVenues(ctx context.Context) {
	venues := vs.VenueRepo.Venues
	for _, venue := range venues {
		if venue.HostID == config.GetUserID(ctx) {
			vs.PrintVenue(venue)
		}
	}
}

func (vs *VenueService) RemoveVenue(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Venue ID to remove: ")
	venueID, _ := reader.ReadString('\n')
	venueID = strings.TrimSpace(venueID)

	venues := vs.VenueRepo.Venues
	for i, venue := range venues {
		if venue.ID == venueID {
			if venue.HostID != config.GetUserID(ctx) {
				fmt.Println("Unauthorized: You cannot remove this venue.")
				return
			}

			vs.PrintVenue(venue)
			fmt.Println("Are you sure you want to delete this venue?")
			fmt.Println("1. Yes")
			fmt.Println("2. No")
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				venues = append(venues[:i], venues[i+1:]...)
				vs.VenueRepo.SaveVenues(venues)
				fmt.Println("Venue removed successfully.")
			}
			return
		}
	}
	fmt.Println("Venue not found.")
}

func (vs *VenueService) PrintVenue(venue models.Venue) {
	fmt.Println("Venue Details:")
	fmt.Println("ID:                    ", venue.ID)
	fmt.Println("Name:                  ", venue.Name)
	fmt.Println("City:                  ", venue.City)
	fmt.Println("State:                 ", venue.State)
	fmt.Println("Seat Layout Required?: ", venue.IsSeatLayoutRequired)
}
