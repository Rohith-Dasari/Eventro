package eventservice

import (
	"bufio"
	"context"
	"eventro2/models"
	eventsrepository "eventro2/repository/event_repository"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type EventService struct {
	EventRepo eventsrepository.EventRepository
}

func NewEventService(eventRepo eventsrepository.EventRepository) *EventService {
	return &EventService{EventRepo: eventRepo}
}

func (e *EventService) ModerateEvents(ctx context.Context) {
	fmt.Println("Enter ID of the event/show to be moderated:")
	var eventID string
	fmt.Scanf("%s", &eventID)

	events := e.EventRepo.Events
	var requiredEvent *models.Event
	var found bool

	for i := range events {
		if events[i].ID == eventID {
			requiredEvent = &events[i]
			e.PrintEvent(*requiredEvent)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Event not found. Please enter a correct ID.")
		return
	}

	if requiredEvent.IsBlocked {
		fmt.Print("Are you sure you want to UNBLOCK the event? (y/n): ")
	} else {
		fmt.Print("Are you sure you want to BLOCK the event? (y/n): ")
	}

	var choice string
	fmt.Scanf("%s", &choice)

	if strings.ToLower(choice) == "y" {
		requiredEvent.IsBlocked = !requiredEvent.IsBlocked
		status := "unblocked"
		if requiredEvent.IsBlocked {
			status = "blocked"
		}
		fmt.Printf("Event successfully %s.\n", status)
	} else {
		fmt.Println("Action canceled.")
	}

	e.EventRepo.SaveEvents(events)
}

func (e *EventService) ViewBlockedEvents(ctx context.Context) {
	events := e.EventRepo.Events
	for _, event := range events {
		if event.IsBlocked {
			e.PrintEvent(event)
		}
	}
}

func (e *EventService) CreateNewEvent() models.Event {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter event name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter event description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Enter number of artists: ")
	var numArtists int
	fmt.Scanf("%d", &numArtists)

	reader.ReadString('\n')
	artists := make([]string, numArtists)
	for i := 0; i < numArtists; i++ {
		fmt.Printf("Enter artist %d: ", i+1)
		artist, _ := reader.ReadString('\n')
		artists[i] = strings.TrimSpace(artist)
	}

	fmt.Print("Enter event duration (e.g., 2h30m): ")
	duration, _ := reader.ReadString('\n')
	duration = strings.TrimSpace(duration)

	var category models.EventCategory
	for {
		fmt.Println("Select category:")
		fmt.Printf("1. %s\n2. %s\n3. %s\n4. %s\n5. %s\n", models.Movie, models.Sports, models.Concert, models.Party, models.Workshop)
		catChoice, _ := reader.ReadString('\n')
		catChoice = strings.TrimSpace(catChoice)
		cat, _ := strconv.Atoi(catChoice)

		switch cat {
		case 1:
			category = models.Movie
		case 2:
			category = models.Sports
		case 3:
			category = models.Concert
		case 4:
			category = models.Party
		case 5:
			category = models.Workshop
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
			continue
		}
		break
	}

	return models.Event{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Artists:     artists,
		Duration:    duration,
		Category:    category,
		IsBlocked:   false,
	}
}

func (e *EventService) PrintEvent(event models.Event) {
	fmt.Println("-------------")
	fmt.Printf("ID: %s\n", event.ID)
	fmt.Printf("Name: %s\n", event.Name)
	fmt.Printf("Description: %s\n", event.Description)
	fmt.Printf("Hype Meter: %d\n", event.HypeMeter)
	fmt.Printf("Duration: %s\n", event.Duration)
	fmt.Printf("Category: %s\n", event.Category)
	if len(event.Artists) > 0 {
		fmt.Printf("Artists: %s\n", strings.Join(event.Artists, ", "))
	}
	fmt.Println("-------------")
}
