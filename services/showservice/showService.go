package showservice

import (
	"context"
	"eventro2/config"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/bookingservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type ShowService struct {
	ShowRepo  showrepository.ShowRepository
	VenueRepo venuerepository.VenueRepository
	BookRepo  bookingrepository.BookingRepository
	EventRepo eventsrepository.EventRepository
}

func NewShowService(
	showRepo showrepository.ShowRepository,
	venueRepo venuerepository.VenueRepository,
	bookRepo bookingrepository.BookingRepository,
	eventRepo eventsrepository.EventRepository,
) ShowService {
	return ShowService{
		ShowRepo:  showRepo,
		VenueRepo: venueRepo,
		BookRepo:  bookRepo,
		EventRepo: eventRepo,
	}
}
func (s *ShowService) BrowseShowsByEvent(ctx context.Context) {
	bookingService := bookingservice.NewBookingService(s.BookRepo, s.ShowRepo, s.VenueRepo, s.EventRepo)

	for {
		fmt.Println("Enter event ID:")
		eventID := utils.ReadLine()

		// fetch shows for this event
		shows, err := s.ShowRepo.ListByEvent(eventID)
		if err != nil {
			color.Red("Failed to fetch shows: %v", err)
			return
		}

		if len(shows) == 0 {
			color.Red("No shows found for this event.")
			fmt.Println("1. Choose another event \n2. Back")
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		found := false
		for _, show := range shows {
			if !show.IsBlocked {
				s.printShow(show)
				found = true
			}
		}
		if !found {
			color.Red("No active shows available for this event.")
			return
		}

		fmt.Print("Enter the date you want to book for (YYYY-MM-DD): ")
		inputDate := utils.ReadLine()
		parsedDate, err := time.Parse("2006-01-02", inputDate)
		if err != nil {
			color.Red("Invalid date format. Please use YYYY-MM-DD.")
			continue
		}

		var availableShows []models.Show
		for _, show := range shows {
			if show.ShowDate.Equal(parsedDate) && !show.IsBlocked {
				s.printShow(show)
				availableShows = append(availableShows, show)
			}
		}

		if len(availableShows) == 0 {
			color.Red("No active shows found for this event on that date.")
			fmt.Println("1. Choose another date \n2. Back")
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		fmt.Println("1. Continue with Show ID\n2. Back")
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		if choice == 1 {
			fmt.Println("Enter show ID:")
			showID := utils.ReadLine()
			s.DisplayShow(ctx, showID)

			var userID string
			if config.GetUserRole(ctx) == models.Admin {
				fmt.Println("Enter the User ID of user you want to book for: ")
				userID = utils.ReadLine()
			} else {
				userID = config.GetUserID(ctx)
			}

			bookingService.MakeBooking(ctx, userID, showID)
		}
		return
	}
}

func (s *ShowService) DisplayShow(ctx context.Context, showID string) {
	// get show by ID
	show, err := s.ShowRepo.GetByID(showID)
	if err != nil || show == nil {
		color.Red("Enter valid Show ID.")
		return
	}

	// fetch venue
	venue, err := s.VenueRepo.GetByID(show.VenueID)
	if err != nil || venue == nil {
		color.Red("Venue not found.")
		return
	}

	// display details
	color.Blue("Show Details")
	fmt.Printf("Show ID:          %s\n", show.ID)
	fmt.Printf("Venue:            %s (%s, %s)\n", venue.Name, venue.City, venue.State)
	fmt.Printf("Price per ticket: ₹%.2f\n", show.Price)

	if venue.IsSeatLayoutRequired {
		fmt.Println("([X] = Booked)")
		s.displaySeatMap(show.BookedSeats)
	} else {
		color.Cyan("This venue does not have a seat layout.")
	}
}

func (s *ShowService) ModerateShow(ctx context.Context) {
	for {
		fmt.Println("Enter ID of show to moderate:")
		showID := utils.ReadLine()

		// fetch show directly
		show, err := s.ShowRepo.GetByID(showID)
		if err != nil || show == nil {
			color.Red("Show not found. Please enter correct Show ID")
			fmt.Println("1. Retry with another Show ID\n2. Back ")
			fmt.Println(config.ChoiceMessage)
			choice, err := utils.TakeUserInput()
			if err != nil {
				color.Red("Please enter an integer")
				continue
			}
			switch choice {
			case 1:
				continue
			default:
				return
			}
		}

		// show details
		s.printShow(*show)

		if show.IsBlocked {
			fmt.Println("The show is currently BLOCKED")
			fmt.Println("Are you sure you want to UNBLOCK the show?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				show.IsBlocked = false
				if err := s.ShowRepo.Update(show); err != nil {
					color.Red("Failed to update show: %v", err)
					return
				}
				fmt.Println("Show is successfully unblocked")
				return
			case 2:
				color.Red("Canceled")
				continue
			default:
				return
			}
		} else {
			fmt.Println("The show is currently UNBLOCKED")
			fmt.Println("Are you sure you want to BLOCK the show?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				show.IsBlocked = true
				if err := s.ShowRepo.Update(show); err != nil {
					color.Red("Failed to update show: %v", err)
					return
				}
				fmt.Println("Show is successfully blocked")
				return
			case 2:
				color.Red("Canceled")
				continue
			default:
				return
			}
		}
	}
}
func (s *ShowService) ViewBlockedShows(ctx context.Context) {
	shows, err := s.ShowRepo.List()
	if err != nil {
		color.Red("Failed to fetch shows: %v", err)
		return
	}

	found := false
	for _, show := range shows {
		if show.IsBlocked {
			s.printShow(show)
			found = true
		}
	}

	if !found {
		color.Red("No Blocked Shows")
	}
}
func (s *ShowService) CreateShow(ctx context.Context) {
	for {
		color.Blue("Let's Create a Show")

		var requiredVenue *models.Venue
		var venueID string
		for {
			fmt.Print("Enter Venue ID: ")
			venueID = utils.ReadLine()

			if venueID == "" {
				color.Red("Venue ID cannot be empty.")
				continue
			}

			venues, err := s.VenueRepo.List()
			if err != nil {
				color.Red("Failed to fetch venues: %v", err)
				return
			}

			for i := range venues {
				if venues[i].ID == venueID {
					requiredVenue = &venues[i]
					break
				}
			}

			if requiredVenue == nil {
				color.Red("Venue ID not found.")
				continue
			}

			if requiredVenue.HostID != config.GetUserID(ctx) {
				color.Red("You are authorised to create shows only on your venues.")
				fmt.Println("1. Enter another Venue ID \n2. Back ")
				choice, _ := utils.TakeUserInput()
				if choice == 1 {
					continue
				} else {
					return
				}
			}
			break
		}

		var eventID string
		for {
			fmt.Print("Enter Event ID: ")
			eventID = strings.TrimSpace(utils.ReadLine())
			if eventID == "" {
				color.Red("Event ID cannot be empty.")
				continue
			}

			event, err := s.EventRepo.GetByID(eventID)
			if err != nil || event == nil {
				color.Red("Event ID not found.")
				continue
			}
			break
		}

		var price float64
		for {
			fmt.Print("Enter Price: ")
			priceInput := utils.ReadLine()
			p, err := strconv.ParseFloat(priceInput, 64)
			if err != nil || p < 0 {
				color.Red("Invalid price. Enter a positive number.")
				continue
			}
			price = p
			break
		}

		var parsedDate time.Time
		for {
			fmt.Print("Enter Show Date (YYYY-MM-DD): ")
			showDate := utils.ReadLine()
			d, err := time.Parse("2006-01-02", showDate)
			if err != nil {
				color.Red("Invalid date format. Please use YYYY-MM-DD.")
				continue
			}
			today := time.Now().Truncate(24 * time.Hour)
			d = d.Truncate(24 * time.Hour)

			if d.Before(today) {
				color.Red("Show date cannot be in the past.")
				continue
			}
			parsedDate = d
			break
		}

		var showTime string
		for {
			fmt.Print("Enter Show Time (HH:MM): ")
			input := utils.ReadLine()
			t, err := time.Parse("15:04", input)
			if err != nil {
				color.Red("Invalid time format. Please use HH:MM in 24-hour format.")
				continue
			}

			today := time.Now()
			if parsedDate.Equal(today.Truncate(24 * time.Hour)) {
				showDateTime := time.Date(
					today.Year(), today.Month(), today.Day(),
					t.Hour(), t.Minute(), 0, 0, time.Local,
				)
				if showDateTime.Before(today) {
					color.Red("Show time must be in the future.")
					continue
				}
			}

			showTime = input
			break
		}
		fmt.Println("1. Confirm \n2. Back")
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		if choice == 2 {
			continue
		}

		show := &models.Show{
			ID:          uuid.New().String(),
			HostID:      config.GetUserID(ctx),
			VenueID:     venueID,
			EventID:     eventID,
			IsBlocked:   false,
			Price:       price,
			ShowDate:    parsedDate,
			ShowTime:    showTime,
			BookedSeats: []string{},
		}

		if err := s.ShowRepo.Create(show); err != nil {
			color.Red("Failed to create show: %v", err)
			return
		}

		color.Green("Show created successfully:")
		s.printShow(*show)
		break
	}
}

func (s *ShowService) BrowseHostShows(ctx context.Context) {
	shows, err := s.ShowRepo.List()
	if err != nil {
		color.Red("Failed to fetch shows: %v", err)
		return
	}

	found := false
	hostID := config.GetUserID(ctx)

	for _, show := range shows {
		if show.HostID == hostID {
			found = true
			s.printShow(show)
			s.displaySeatMap(show.BookedSeats)
		}
	}

	if !found {
		color.Red("No shows found for your host account.")
	}
}
func (s *ShowService) RemoveHostShow(ctx context.Context) {
	for {
		fmt.Print("Enter Show ID to remove: ")
		showID := strings.TrimSpace(utils.ReadLine())

		if showID == "" {
			color.Red("Show ID cannot be empty.")
			continue
		}

		show, err := s.ShowRepo.GetByID(showID)
		if err != nil || show == nil {
			color.Red("Show not found.")
			fmt.Println("1. Retry \n2. Back")
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		if show.HostID != config.GetUserID(ctx) {
			color.Red("You are authorised to delete only your shows")
			fmt.Println("1. Retry \n2. Back")
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			return
		}

		s.printShow(*show)

		fmt.Println("Are you sure you want to delete this show?")
		fmt.Println("1. Yes")
		fmt.Println("2. No, another ID")
		fmt.Println("3. Back")

		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			if err := s.ShowRepo.Delete(showID); err != nil {
				color.Red("Failed to delete show: %v", err)
				return
			}
			color.Green("Show removed successfully.")
			return
		case 2:
			color.Yellow("Canceled. Try another ID.")
			continue
		case 3:
			return
		default:
			color.Red("Invalid choice.")
		}
	}
}

func (s *ShowService) printShow(show models.Show) {
	venue, _ := s.VenueRepo.GetByID(show.VenueID)

	fmt.Println(config.Dash)
	fmt.Printf("%-15s %s\n", "Show ID:", show.ID)
	fmt.Printf("%-15s %s\n", "Show Date:", show.ShowDate.Format("2006-01-02"))
	fmt.Printf("%-15s %s\n", "Show Time:", show.ShowTime) // already string
	if venue != nil {
		fmt.Printf("%-15s %s\n", "Venue Name:", venue.Name)
		fmt.Printf("%-15s %s\n", "Venue City:", venue.City)
	}
	fmt.Printf("%-15s ₹%.2f\n", "Price:", show.Price)
	fmt.Println(config.Dash)
}

func (s *ShowService) displaySeatMap(booked []string) {
	rows := config.Rows
	cols := config.Columns

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
