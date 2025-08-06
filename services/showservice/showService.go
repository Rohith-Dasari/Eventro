package showservice

import (
	"context"
	"eventro/config"
	"eventro/models"
	"eventro/storage"
	"fmt"
	"strings"
)

func BrowseShowsByEvent(eventID string, shows []models.Show) {
	//unmarshall shows file into slices of shows and display shows which has the eventID same as selected
	found := false
	fmt.Printf("Available Shows for the Event:\n")
	for _, show := range shows {
		if show.EventID == eventID && !show.IsBlocked {
			printShow(show)
			//create map of show and venue(name, city )
			found = true
		}
	}
	if !found {
		fmt.Println("No active shows found for this event.")
	}
}

func DisplayShow(showID string, shows []models.Show) {
	var selectedShow *models.Show

	for i := range shows {
		if shows[i].ID == showID {
			selectedShow = &shows[i]
			break
		}
	}

	if selectedShow == nil {
		fmt.Println("Show not found")
		return
	}

	venues := storage.LoadVenues()
	var venue *models.Venue
	for i := range venues {
		if venues[i].ID == selectedShow.VenueID {
			venue = &venues[i]
			break
		}
	}
	if venue == nil {
		fmt.Println("Venue not found.")
		return
	}

	fmt.Println("Show Details")
	fmt.Printf("Show ID: %s\n", selectedShow.ID)
	fmt.Printf("Venue: %s (%s, %s)\n", venue.Name, venue.City, venue.State)
	fmt.Printf("Price per ticket: ₹%.2f\n", selectedShow.Price)

	if venue.IsSeatLayoutRequired {
		fmt.Println("([X] = Booked)")
		displaySeatMap(selectedShow.BookedSeats)
	} else {
		fmt.Println("This venue does not have a seat layout.")
	}
}

func printShow(s models.Show) {
	//venue as well-
	fmt.Printf("Show ID: %s\n", s.ID)
	fmt.Printf("Show Time: %s\n", s.ShowTime)
	fmt.Printf("Price: ₹%.2f\n", s.Price)
	fmt.Println("-------------------------")
}

func displaySeatMap(booked []string) {
	const rows = 10
	const cols = 10

	bookedMap := make(map[string]bool)
	for _, seat := range booked {
		bookedMap[strings.ToUpper(seat)] = true
	}

	for i := 0; i < rows; i++ {
		rowLabel := "A" + string(i)
		fmt.Printf("%s  ", rowLabel)
		for j := 1; j <= cols; j++ {
			seat := fmt.Sprintf("%s%d", rowLabel, j)
			if bookedMap[seat] {
				fmt.Print(config.BookedSeat)
			} else {
				fmt.Printf(config.AvailableSeat, seat)
			}
		}
		fmt.Println()
	}
}

func BlockShow(ctx context.Context) {
	fmt.Println("enter ID of show want to block")
	var showID string
	fmt.Scanf("%s", &showID)
	shows := storage.LoadShows()
	var requiredShow *models.Show
	var found bool
	for _, show := range shows {
		if show.ID == showID {
			requiredShow = &show
			printShow(show)
			found = true
		}
	}
	if !found {
		fmt.Println("Show not found, please enter correct ID")
	} else {
		fmt.Printf("Are you sure you want to block the show: y/n")
		var choice string
		fmt.Scanf("%s", choice)
		if choice == "y" {
			requiredShow.IsBlocked = true
		}
	}
}
