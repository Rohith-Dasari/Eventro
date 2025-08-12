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
}

func NewShowService(showRepo showrepository.ShowRepository, venueRepo venuerepository.VenueRepository) *ShowService {
	return &ShowService{
		ShowRepo:  showRepo,
		VenueRepo: venueRepo,
	}
}
func (s *ShowService) BrowseShowsByEvent(ctx context.Context) {
	showRepo := showrepository.NewShowRepository()
	bookingRepo := bookingrepository.NewBoookingStore()
	bookingService := bookingservice.NewBookingService(*bookingRepo, *showRepo)
	shows := s.ShowRepo.Shows
	// reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter event ID:")
		eventID := utils.ReadLine()
		found := false
		for _, show := range shows {
			if show.EventID == eventID && !show.IsBlocked {
				s.printShow(show)
				found = true
			}
		}
		if !found {
			color.Red("No active shows found for this event.")
			fmt.Println("1. Choose another event \n2. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			default:
				return
			}
		}

		fmt.Print("Enter the date you want to book for (YYYY-MM-DD): ")
		inputDate := utils.ReadLine()
		found = false
		for _, show := range shows {
			if show.EventID == eventID && show.ShowDate == inputDate && !show.IsBlocked {
				s.printShow(show)
				found = true
			}
		}
		if !found {
			color.Red("No active shows found for this event.")
			fmt.Println("1. Choose another date\n2. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			default:
				return
			}
		} else {
			fmt.Println("1. Continue with Show ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
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
			default:
				continue
			}
		}
		return
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
		color.Red("Enter valid Show ID.")
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
		color.Red("Venue not found.")
		return
	}

	color.Blue("Show Details")
	fmt.Printf("Show ID:          %s\n", selectedShow.ID)
	fmt.Printf("Venue:            %s (%s, %s)\n", venue.Name, venue.City, venue.State)
	fmt.Printf("Price per ticket: ₹%.2f\n", selectedShow.Price)
	if venue.IsSeatLayoutRequired {
		fmt.Println("([X] = Booked)")
		s.displaySeatMap(selectedShow.BookedSeats)
	} else {
		color.Cyan("This venue does not have a seat layout.")
	}
}

func (s *ShowService) ModerateShow(ctx context.Context) {
	for {
		fmt.Println("Enter ID of show to moderate:")
		showID := utils.ReadLine()

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
			color.Red("Show not found. Please enter correct Show ID")
			fmt.Println("1. Retry with another Show ID\n2. Back ")
			fmt.Println(config.ChoiceMessage)
			choice, err := utils.TakeUserInput()
			if err != nil {
				fmt.Println("please enter an integer")
				continue
			}

			switch choice {
			case 1:
				continue
			case 2:
				return
			default:
				return
			}
		}
		if requiredShow.IsBlocked {
			fmt.Println("The show is currently BLOCKED")
			fmt.Println("Are you sure you want to UNBLOCK the show?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				requiredShow.IsBlocked = !requiredShow.IsBlocked
				s.ShowRepo.SaveShows(shows)
				fmt.Println("Show is successfully unblocked")
				return
			case 2:
				color.Red("Canceled")
				continue
			case 3:
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
				requiredShow.IsBlocked = !requiredShow.IsBlocked
				s.ShowRepo.SaveShows(shows)
				fmt.Println("Show is successfully blocked")
				return
			case 2:
				color.Red("Canceled")
				continue
			case 3:
				return
			}
		}
	}
}

func (s *ShowService) ViewBlockedShows(ctx context.Context) {
	shows := s.ShowRepo.Shows
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
		color.Blue("Lets Create a Show")

		var requiredVenue models.Venue
		var venueID string
		for {
			fmt.Print("Enter Venue ID: ")
			venueID = utils.ReadLine()

			if venueID == "" {
				color.Red("Venue ID cannot be empty.")
				continue
			}

			found := false
			for _, venue := range s.VenueRepo.Venues {
				if venue.ID == venueID {
					requiredVenue = venue
					found = true
					break
				}
			}
			if !found {
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

		fmt.Print("Enter Event ID: ")
		var eventID string
		for {
			fmt.Print("Enter Event ID: ")
			eventID = strings.TrimSpace(utils.ReadLine())
			if eventID == "" {
				color.Red("Event ID cannot be empty.")
				continue
			}
			eventRepo := eventsrepository.NewEventRepository()
			found := false
			for _, event := range eventRepo.Events {
				if event.ID == eventID {
					found = true
					break
				}
			}
			if !found {
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

		var showDate string
		for {
			fmt.Print("Enter Show Date (YYYY-MM-DD): ")
			showDate = utils.ReadLine()
			parsedDate, err := time.Parse("2006-01-02", showDate)
			if err != nil {
				color.Red("Invalid date format. Please use YYYY-MM-DD.")
				continue
			}
			today := time.Now().Truncate(24 * time.Hour)
			parsedDate = parsedDate.Truncate(24 * time.Hour)

			if parsedDate.Before(today) {
				color.Red("Show date cannot be in the past.")
				continue
			}

			break
		}

		var showTime string
		for {
			fmt.Print("Enter Show Time (HH:MM): ")
			showTime = utils.ReadLine()
			_, err := time.Parse("15:04", showTime)
			if err != nil {
				color.Red("Invalid time format. Please use HH:MM in 24-hour format.")
				continue
			}
			break
		}

		fmt.Println("1. Confirm \n2. Back")
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		if choice == 2 {
			continue
		}

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
		s.ShowRepo.Shows = append(s.ShowRepo.Shows, show)

		fmt.Println("Show created successfully:")
		s.printShow(show)
		break
	}
}

func (s *ShowService) BrowseHostShows(ctx context.Context) {
	shows := s.ShowRepo.Shows
	found := false
	for _, show := range shows {
		if show.HostID == config.GetUserID(ctx) {
			found = true
			s.printShow(show)
			s.displaySeatMap(show.BookedSeats)
		}
	}
	if !found {
		fmt.Println("Show not found")
	}
}

func (s *ShowService) RemoveHostShow(ctx context.Context) {
	for {
		fmt.Print("Enter Show ID to remove: ")
		showID := utils.ReadLine()

		shows := s.ShowRepo.Shows
		var requiredShow models.Show
		var requiredIndex int
		var found bool
		for i, show := range shows {
			if show.ID == showID {
				requiredShow = show
				requiredIndex = i
				found = true
				break
			}
		}

		if !found || requiredShow.HostID != config.GetUserID(ctx) {
			color.Red("You are authorised to delete only your shows")

			fmt.Println("1. Retry \n2. Back")
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

			s.printShow(requiredShow)
		loop:
			for {
				fmt.Println("Are you sure you want to delete this show?")
				fmt.Println("1. Yes")
				fmt.Println("2. No, another ID")
				fmt.Println("3. Back")
				choice, _ := utils.TakeUserInput()
				if choice == 1 {
					s.ShowRepo.Shows = append(s.ShowRepo.Shows[:requiredIndex], s.ShowRepo.Shows[requiredIndex+1:]...)
					s.ShowRepo.SaveShows(s.ShowRepo.Shows)
					fmt.Println("Show removed successfully.")
					return
				} else if choice == 2 {
					color.Red("Canceled")
					break loop
				} else {
					return
				}
			}
		}
	}
}

func (s *ShowService) printShow(show models.Show) {
	var requiredVenue models.Venue
	for _, venue := range s.VenueRepo.Venues {
		if venue.ID == show.VenueID {
			requiredVenue = venue
		}
	}
	fmt.Println(config.Dash)
	fmt.Printf("%-15s %s\n", "Show ID:", show.ID)
	fmt.Printf("%-15s %s\n", "Show Date:", show.ShowDate)
	fmt.Printf("%-15s %s\n", "Show Time:", show.ShowTime)
	fmt.Printf("%-15s %s\n", "Venue Name:", requiredVenue.Name)
	fmt.Printf("%-15s %s\n", "Venue City:", requiredVenue.City)
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
