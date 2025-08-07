package venueservice

import (
	"bufio"
	"context"
	"eventro/config"
	"eventro/models"
	"eventro/storage"
	utils "eventro/utils/userinput"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func AddVenue(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Venue Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Enter Host ID: ")
	hostID, _ := reader.ReadString('\n')

	fmt.Print("Enter City: ")
	city, _ := reader.ReadString('\n')

	fmt.Print("Enter State: ")
	state, _ := reader.ReadString('\n')

	fmt.Print("Is seat layout required? (y/n): ")
	seatLayoutInput, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	hostID = strings.TrimSpace(hostID)
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)
	seatLayoutInput = strings.TrimSpace(strings.ToLower(seatLayoutInput))
	var isSeatLayoutRequired bool
	if seatLayoutInput == "y" {
		isSeatLayoutRequired = true
	} else {
		isSeatLayoutRequired = false
	}
	venue := models.Venue{
		ID:                   uuid.New().String(),
		Name:                 name,
		HostID:               hostID,
		City:                 city,
		State:                state,
		IsSeatLayoutRequired: isSeatLayoutRequired,
	}
	venues := storage.LoadVenues()
	venues = append(venues, venue)
	storage.SaveVenues(venues)
	fmt.Println("Venue created successfully:")
	fmt.Printf("%+v\n", venue)
}
func BrowseHostVenues(ctx context.Context) {
	venues := storage.LoadVenues()
	for _, venue := range venues {
		if venue.HostID == config.GetUserID(ctx) {
			PrintVenue(venue)
		}
	}

}

func RemoveVenue(ctx context.Context) {
	fmt.Println("Enter Venue id need to be removed")
	reader := bufio.NewReader(os.Stdin)
	venueID, _ := reader.ReadString('\n')
	venueID = strings.TrimSpace(venueID)
	venues := storage.LoadVenues()

	for index, venue := range venues {
		if venue.ID == venueID {
			if venue.HostID == config.GetUserID(ctx) {
				PrintVenue(venue)
				fmt.Println("Are you sure you want to delete this ")
				fmt.Println("1. yes 2. no, go back")
				choice, _ := utils.TakeUserInput()
				if choice == 1 {
					venues = append(venues[:index], venues[index+1:]...)
					storage.SaveVenues(venues)
				}
			} else {
				fmt.Println("You dont have access to remove the venue.")

			}
		} else {
			fmt.Println("Venue not found")
		}
	}
}

func PrintVenue(venue models.Venue) {
	fmt.Println("Venue Details:")
	fmt.Println("ID:                    ", venue.ID)
	fmt.Println("Name:                  ", venue.Name)
	fmt.Println("City:                  ", venue.City)
	fmt.Println("State:                 ", venue.State)
	fmt.Println("Seat Layout Required?: ", venue.IsSeatLayoutRequired)
}
