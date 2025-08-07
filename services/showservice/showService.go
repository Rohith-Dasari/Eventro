package showservice

import (
	"bufio"
	"context"
	"eventro/config"
	"eventro/models"
	"eventro/storage"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func BrowseShowsByEvent(eventID string, shows []models.Show) {
	//unmarshall shows file into slices of shows and display shows which has the eventID and date same as selected
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the date you want to book for (YYYY-MM-DD): ")
	inputDate, _ := reader.ReadString('\n')
	inputDate = strings.TrimSpace(inputDate)
	found := false
	fmt.Printf("Available Shows for the Event:\n")
	for _, show := range shows {
		if show.EventID == eventID && show.ShowDate == inputDate && !show.IsBlocked {
			printShow(show)
			//create map of show and venue(name, city )
			//
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
		fmt.Println("How many tickets do you want to book")
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
		rowLabel := string('A' + i)
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

func ModerateShow(ctx context.Context) {
	fmt.Println("enter id of show to be blocked")
	var showID string
	fmt.Scanf("%s", &showID)
	shows := storage.LoadShows()
	var requiredShow *models.Show
	var found bool
	for _, show := range shows {
		if show.ID == showID {
			requiredShow = &show
			printShow(*requiredShow)
			found = true
		}
	}
	if !found {
		fmt.Println("show not found, please enter correct ID")
	} else {
		if !requiredShow.IsBlocked {
			fmt.Print("Are you sure you want to unblock the show: y/n")
			requiredShow.IsBlocked = true
		} else {
			fmt.Printf("Are you sure you want to block the user: y/n")
			var choice string
			fmt.Scanf("%s", choice)
			if choice == "y" {
				requiredShow.IsBlocked = true
			}
		}
	}
	storage.SaveShows(shows)
}

func ViewBlockedShows(ctx context.Context) {
	shows := storage.LoadShows()
	for _, show := range shows {
		printShow(show)
	}
}

func CreateShow(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Venue ID: ")
	venueID, _ := reader.ReadString('\n')

	fmt.Print("Enter Event ID: ")
	eventID, _ := reader.ReadString('\n')

	fmt.Print("Enter Price: ")
	var price float64
	fmt.Scanf("%f\n", &price)

	fmt.Print("Enter Show Date (YYYY-MM-DD): ")
	showDate, _ := reader.ReadString('\n')

	fmt.Print("Enter Show Time (HH:MM): ")
	showTime, _ := reader.ReadString('\n')

	venueID = strings.TrimSpace(venueID)
	eventID = strings.TrimSpace(eventID)
	showDate = strings.TrimSpace(showDate)
	showTime = strings.TrimSpace(showTime)

	show := models.Show{
		ID:          uuid.New().String(),
		HostID:      config.GetUserID(ctx),
		VenueID:     venueID,
		EventID:     eventID,
		CreatedAt:   time.Now().Format(time.RFC3339),
		IsBlocked:   false,
		Price:       price,
		ShowDate:    showDate,
		ShowTime:    showTime,
		BookedSeats: []string{},
	}

	shows := storage.LoadShows()
	shows = append(shows, show)
	storage.SaveShows(shows)

	fmt.Println("Show created successfully:")
	fmt.Printf("%+v\n", show)
}
