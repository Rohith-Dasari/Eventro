package showservice

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ShowService struct {
	ShowRepo  showrepository.ShowRepository
	VenueRepo venuerepository.VenueRepository
}

func NewShowService(showRepo showrepository.ShowRepository, venueRepo venuerepository.VenueRepository) *ShowService {
	return &ShowService{
		ShowRepo:  showRepo,
		VenueRepo: venueRepo,
	}
}

func (s *ShowService) BrowseShowsByEvent(ctx context.Context, eventID string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the date you want to book for (YYYY-MM-DD): ")
	inputDate, _ := reader.ReadString('\n')
	inputDate = strings.TrimSpace(inputDate)

	shows := s.ShowRepo.Shows
	found := false
	for _, show := range shows {
		if show.EventID == eventID && show.ShowDate == inputDate && !show.IsBlocked {
			s.printShow(show)
			found = true
		}
	}
	if !found {
		fmt.Println("No active shows found for this event.")
	}
}

func (s *ShowService) DisplayShow(ctx context.Context, showID string) {
	shows := s.ShowRepo.Shows
	var selectedShow *models.Show
	for i := range shows {
		if shows[i].ID == showID {
			selectedShow = &shows[i]
			break
		}
	}

	if selectedShow == nil {
		fmt.Println("Show not found.")
		return
	}

	venues := s.VenueRepo.Venues
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
		s.displaySeatMap(selectedShow.BookedSeats)
	} else {
		fmt.Println("This venue does not have a seat layout.")
	}
}

func (s *ShowService) ModerateShow(ctx context.Context) {
	fmt.Println("Enter ID of show to moderate:")
	var showID string
	fmt.Scanf("%s", &showID)

	shows := s.ShowRepo.Shows
	var requiredShow *models.Show
	var found bool
	for i := range shows {
		if shows[i].ID == showID {
			requiredShow = &shows[i]
			s.printShow(*requiredShow)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Show not found.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	if requiredShow.IsBlocked {
		fmt.Print("Unblock this show? (y/n): ")
	} else {
		fmt.Print("Block this show? (y/n): ")
	}
	choice, _ := reader.ReadString('\n')
	if strings.TrimSpace(choice) == "y" {
		requiredShow.IsBlocked = !requiredShow.IsBlocked
		fmt.Println("Show moderation status changed.")
		s.ShowRepo.SaveShows(shows)
	} else {
		fmt.Println("No changes made.")
	}
}

func (s *ShowService) ViewBlockedShows(ctx context.Context) {
	shows := s.ShowRepo.Shows
	for _, show := range shows {
		if show.IsBlocked {
			s.printShow(show)
		}
	}
}

func (s *ShowService) CreateShow(ctx context.Context) {
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

	show := models.Show{
		ID:          uuid.New().String(),
		HostID:      config.GetUserID(ctx),
		VenueID:     strings.TrimSpace(venueID),
		EventID:     strings.TrimSpace(eventID),
		CreatedAt:   time.Now().Format(time.RFC3339),
		IsBlocked:   false,
		Price:       price,
		ShowDate:    strings.TrimSpace(showDate),
		ShowTime:    strings.TrimSpace(showTime),
		BookedSeats: []string{},
	}

	shows := s.ShowRepo.Shows
	shows = append(shows, show)
	s.ShowRepo.SaveShows(shows)

	fmt.Println("Show created successfully:")
	s.printShow(show)
}

func (s *ShowService) BrowseHostShows(ctx context.Context) {
	shows := s.ShowRepo.Shows
	for _, show := range shows {
		if show.HostID == config.GetUserID(ctx) {
			s.printShow(show)
			s.displaySeatMap(show.BookedSeats)
		}
	}
}

func (s *ShowService) RemoveHostShow(ctx context.Context) {
	fmt.Print("Enter Show ID to remove: ")
	reader := bufio.NewReader(os.Stdin)
	showID, _ := reader.ReadString('\n')
	showID = strings.TrimSpace(showID)

	shows := s.ShowRepo.Shows
	for i, show := range shows {
		if show.ID == showID {
			if show.HostID != config.GetUserID(ctx) {
				fmt.Println("Unauthorized: You cannot delete this show.")
				return
			}

			s.printShow(show)
			fmt.Println("Are you sure you want to delete this show?")
			fmt.Println("1. Yes")
			fmt.Println("2. No")
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				shows = append(shows[:i], shows[i+1:]...)
				s.ShowRepo.SaveShows(shows)
				fmt.Println("Show removed successfully.")
			}
			return
		}
	}
	fmt.Println("Show not found.")
}

func (s *ShowService) printShow(show models.Show) {
	fmt.Println("-------------------------")
	fmt.Printf("Show ID: %s\n", show.ID)
	fmt.Printf("Show Date: %s\n", show.ShowDate)
	fmt.Printf("Show Time: %s\n", show.ShowTime)
	fmt.Printf("Price: ₹%.2f\n", show.Price)
	fmt.Println("-------------------------")
}

func (s *ShowService) displaySeatMap(booked []string) {
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
