package searchevents

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	artistrepository "eventro2/repository/artists_repository"
	bookingrepository "eventro2/repository/booking_repository"
	eventrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/artistservice"
	"eventro2/services/showservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type SearchService struct {
	EventRepo  eventrepository.EventRepository
	VenueRepo  venuerepository.VenueRepository
	ShowRepo   showrepository.ShowRepository
	BookRepo   bookingrepository.BookingRepository
	ArtistRepo artistrepository.ArtistRepository
}

func NewSearchService(
	eventRepo eventrepository.EventRepository,
	venueRepo venuerepository.VenueRepository,
	showRepo showrepository.ShowRepository,
	bookRepo bookingrepository.BookingRepository,
	artistRepo artistrepository.ArtistRepository,
) SearchService {
	return SearchService{
		EventRepo:  eventRepo,
		VenueRepo:  venueRepo,
		ShowRepo:   showRepo,
		BookRepo:   bookRepo,
		ArtistRepo: artistRepo,
	}
}

func (s *SearchService) Search(ctx context.Context) {

	for {
		showService := showservice.NewShowService(s.ShowRepo, s.VenueRepo, s.BookRepo, s.EventRepo)
		fmt.Println(config.SearchEventsMessage)
		fmt.Println("1. Search by Event Name")
		fmt.Println("2. Search by Category")
		fmt.Println("3. Search by Location")
		fmt.Println("4. Search by Artist")
		fmt.Println("5. Back")
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()

		switch choice {
		case 1:
			s.SearchByEventName()
			fmt.Println("1. Continue with Event ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				showService.BrowseShowsByEvent(ctx)
			default:
				continue
			}
		case 2:
			s.SearchByCategory()
			fmt.Println("1. Continue with Event ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				showService.BrowseShowsByEvent(ctx)
			default:
				continue
			}
		case 3:
			s.SearchByLocation()
			fmt.Println(config.ChoiceMessage)
			fmt.Println("1. Continue with Event ID\n2. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				showService.BrowseShowsByEvent(ctx)
			default:
				continue
			}
		case 4:
			as := artistservice.NewArtistService(s.ArtistRepo)
			as.BrowseArtist(ctx)
			fmt.Println("1. Continue with Event ID\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				showService.BrowseShowsByEvent(ctx)
			default:
				continue
			}

		case 5:
			return
		default:
			fmt.Println(config.InvalidMSG)
			continue
		}
	}
}

func (s *SearchService) SearchByEventName() {
	for {
		fmt.Print("Enter the name of the event: ")
		query := utils.ReadLine()
		query = strings.ToLower(strings.TrimSpace(query))

		var found bool
		fmt.Println("Matching Events:")

		events, err := s.EventRepo.List()
		if err != nil {
			color.Red("Failed to fetch events: %v", err)
			return
		}

		for _, event := range events {
			if strings.Contains(strings.ToLower(event.Name), query) {
				s.printEvent(event)
				found = true
			}
		}

		if !found {
			color.Red("No matching events found.")
			fmt.Println("1. Enter another event name \n2. Back")
			fmt.Print(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
		}
		break
	}
}

func (s *SearchService) SearchByCategory() {
	for {
		fmt.Printf("Available Categories: \n1. %s \n2. %s \n3. %s \n4. %s \n5. %s\n",
			models.Movie, models.Sports, models.Concert, models.Workshop, models.Party)

		fmt.Print(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()

		var query models.EventCategory
		switch choice {
		case 1:
			query = models.Movie
		case 2:
			query = models.Sports
		case 3:
			query = models.Concert
		case 4:
			query = models.Workshop
		case 5:
			query = models.Party
		default:
			fmt.Println(config.DefaultChoiceMessage)
			continue
		}

		var found bool
		fmt.Println("Events in selected category:")

		events, err := s.EventRepo.List()
		if err != nil {
			color.Red("Failed to fetch events: %v", err)
			return
		}

		for _, event := range events {
			if event.Category == query {
				s.printEvent(event)
				found = true
			}
		}

		if !found {
			color.Red("No events found in this category.")
			fmt.Println("1. Enter another category\n2. Back")
			fmt.Print(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
		}
		break
	}
}

func (s *SearchService) SearchByLocation() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter city to search for events: ")
		input, _ := reader.ReadString('\n')
		city := strings.TrimSpace(input)
		if city == "" {
			color.Red("City cannot be empty.")
			continue
		}

		events, err := s.EventRepo.GetEventsByCity(city)
		if err != nil {
			color.Red("Failed to fetch events: %v", err)
			return
		}

		if len(events) == 0 {
			color.Red("No events found in %s.", city)
			fmt.Println("1. Enter another city \n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			if choice == 1 {
				continue
			}
			break
		}

		for _, event := range events {
			if event.IsBlocked {
				continue
			}
			s.printEvent(event)
		}
		break
	}
}

// func contains(cities []string, target string) bool {
// 	for _, city := range cities {
// 		if strings.EqualFold(city, target) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (s *SearchService) printEvent(event models.Event) {
	fmt.Println(config.Dash)
	fmt.Printf("%-12s : %s\n", "ID", event.ID)
	fmt.Printf("%-12s : %s\n", "Name", event.Name)
	fmt.Printf("%-12s : %s\n", "Description", event.Description)
	fmt.Printf("%-12s : %d\n", "Hype Meter", event.HypeMeter)
	fmt.Printf("%-12s : %s\n", "Duration", event.Duration)
	fmt.Printf("%-12s : %s\n", "Category", event.Category)

	artists, err := s.EventRepo.GetArtistsByEventID(event.ID)
	if err != nil {
		fmt.Println("Error fetching artists:", err)
		return
	}
	fmt.Println("Artists:")
	for _, artist := range artists {
		fmt.Printf("%s,\t", artist.Name)
	}

	fmt.Println(config.Dash)
}
