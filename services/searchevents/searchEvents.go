package searchevents

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	eventrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/showservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type SearchService struct {
	EventRepo eventrepository.EventStorageI
}

func NewSearchService(repo eventrepository.EventStorageI) *SearchService {
	return &SearchService{
		EventRepo: repo,
	}
}

func (s *SearchService) Search(ctx context.Context) {

	for {
		showRepo := showrepository.NewShowRepository()
		venueRepo := venuerepository.NewVenueRepository()
		showService := showservice.NewShowService(*showRepo, *venueRepo)
		fmt.Println(config.SearchEventsMessage)
		fmt.Println("1. Search by Event Name")
		fmt.Println("2. Search by Category")
		fmt.Println("3. Search by Location")
		fmt.Println("4. Back")
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

		var found bool
		fmt.Println("Matching Events:")
		events, _ := s.EventRepo.GetEvents()
		for _, event := range events {
			if strings.Contains(strings.ToLower(event.Name), strings.ToLower(query)) {
				s.printEvent(event)
				found = true
			}
		}
		if !found {
			color.Red("No matching events found.")
			fmt.Println("1. Enter another event name 2. Back")
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			default:
				break
			}
		}
		break
	}

}

func (s *SearchService) SearchByCategory() {
	for {
		fmt.Printf("Available Categories: \n1. %s \n2. %s \n3. %s \n4. %s \n5. %s", models.Movie, models.Sports, models.Concert, models.Workshop, models.Party)
		fmt.Print(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		var query string
		switch choice {
		case 1:
			query = string(models.Movie)
		case 2:
			query = string(models.Sports)
		case 3:
			query = string(models.Concert)
		case 4:
			query = string(models.Workshop)
		case 5:
			query = string(models.Party)
		default:
			fmt.Println(config.DefaultChoiceMessage)
		}

		var found bool
		fmt.Println("Events in selected category:")
		events, _ := s.EventRepo.GetEvents()
		for _, event := range events {
			if strings.ToLower(string(event.Category)) == query {
				s.printEvent(event)
				found = true
			}
		}
		if !found {
			color.Red("No events found in this category.")
			fmt.Println("1. Enter another category\n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			default:
				break
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
		city := strings.ToLower(strings.TrimSpace(input))

		var found bool
		fmt.Printf("\nEvents in %s:\n", city)
		events, _ := s.EventRepo.GetEvents()
		for _, event := range events {
			if event.IsBlocked {
				continue
			}
			if contains(event.Locations, city) {
				s.printEvent(event)
				found = true
			}
		}
		if !found {
			color.Red("No events found in this city.")
			fmt.Println("1. Enter another city \n2. Back")
			fmt.Println(config.ChoiceMessage)
			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				continue
			default:
				break
			}
		}
		break
	}
}

func contains(cities []string, target string) bool {
	for _, city := range cities {
		if strings.EqualFold(city, target) {
			return true
		}
	}
	return false
}

func (s *SearchService) printEvent(e models.Event) {
	fmt.Println(config.Dash)
	fmt.Printf("%-15s %s\n", "ID:", e.ID)
	fmt.Printf("%-15s %s\n", "Name:", e.Name)
	fmt.Printf("%-15s %s\n", "Description:", e.Description)
	fmt.Printf("%-15s %d\n", "Hype Meter:", e.HypeMeter)
	fmt.Printf("%-15s %s\n", "Duration:", e.Duration)
	fmt.Printf("%-15s %s\n", "Category:", e.Category)
	if len(e.Artists) > 0 {
		fmt.Printf("%-15s %s\n", "Artists:", strings.Join(e.Artists, ", "))
	}
	fmt.Println(config.Dash)
}
