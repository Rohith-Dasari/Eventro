package searchevents

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	eventrepository "eventro2/repository/event_repository"
	"fmt"
	"os"
	"strings"
)

type SearchService struct {
	EventRepo eventrepository.EventRepository
}

func NewSearchService(repo eventrepository.EventRepository) *SearchService {
	return &SearchService{
		EventRepo: repo,
	}
}

func (s *SearchService) Search(ctx context.Context) {
	fmt.Println("Select how you want to search")
	fmt.Println("1. Search by Event Name")
	fmt.Println("2. Search by Category")
	fmt.Println("3. Search by Location")
	fmt.Println("4. Back")
	var choice int
	fmt.Print(config.ChoiceMessage)
	fmt.Scanf("%d\n", &choice)

	switch choice {
	case 1:
		s.SearchByEventName()
	case 2:
		s.SearchByCategory()
	case 3:
		s.SearchByLocation()
	case 4:
		return
	default:
		fmt.Println("Invalid choice.")
	}
}

func (s *SearchService) SearchByEventName() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the event: ")
	input, _ := reader.ReadString('\n')
	query := strings.TrimSpace(input)

	var found bool
	fmt.Println("Matching Events:")
	for _, event := range s.EventRepo.Events {
		if strings.Contains(strings.ToLower(event.Name), strings.ToLower(query)) {
			s.printEvent(event)
			found = true
		}
	}
	if !found {
		fmt.Println("No matching events found.")
	}
}

func (s *SearchService) SearchByCategory() {
	fmt.Println("Available Categories: movie, sports, concert, workshop, party")
	fmt.Print("Enter category: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	query := strings.ToLower(strings.TrimSpace(input))

	var found bool
	fmt.Println("Events in selected category:")
	for _, event := range s.EventRepo.Events {
		if strings.ToLower(string(event.Category)) == query {
			s.printEvent(event)
			found = true
		}
	}
	if !found {
		fmt.Println("No events found in this category.")
	}
}

func (s *SearchService) SearchByLocation() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter city to search for events: ")
	input, _ := reader.ReadString('\n')
	city := strings.ToLower(strings.TrimSpace(input))

	var found bool
	fmt.Printf("\nEvents in %s:\n", city)
	for _, event := range s.EventRepo.Events {
		if event.IsBlocked {
			continue
		}
		if contains(event.Locations, city) {
			s.printEvent(event)
			found = true
		}
	}
	if !found {
		fmt.Println("No events found in this city.")
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
	fmt.Println("-------------")
	fmt.Printf("ID: %s\n", e.ID)
	fmt.Printf("Name: %s\n", e.Name)
	fmt.Printf("Description: %s\n", e.Description)
	fmt.Printf("Hype Meter: %d\n", e.HypeMeter)
	fmt.Printf("Duration: %s\n", e.Duration)
	fmt.Printf("Category: %s\n", e.Category)
	if len(e.Artists) > 0 {
		fmt.Printf("Artists: %s\n", strings.Join(e.Artists, ", "))
	}
	fmt.Println("-------------")
}
