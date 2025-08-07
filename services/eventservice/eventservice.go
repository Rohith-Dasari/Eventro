package eventservice

import (
	"bufio"
	"context"
	"eventro/models"
	"eventro/storage"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ModerateEvents(ctx context.Context) {
	fmt.Println("enter id of event of show to be blocked")
	var eventID string
	fmt.Scanf("%s", &eventID)
	events := storage.LoadEvents()
	var requiredEvent *models.Event
	var found bool
	for _, event := range events {
		if event.ID == eventID {
			requiredEvent = &event
			PrintEvent(*requiredEvent)
			found = true
		}
	}
	if !found {
		fmt.Println("event not found, please enter correct ID")
	} else {
		if !requiredEvent.IsBlocked {
			fmt.Print("Are you sure you want to unblock the show: y/n")
			requiredEvent.IsBlocked = true
		} else {
			fmt.Printf("Are you sure you want to block the user: y/n")
			var choice string
			fmt.Scanf("%s", choice)
			if choice == "y" {
				requiredEvent.IsBlocked = true
			}
		}
	}
	storage.SaveEvents(events)
}

func PrintEvent(e models.Event) {
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

func ViewBlockedEvents(ctx context.Context) {
	events := storage.LoadEvents()
	for _, event := range events {
		if event.IsBlocked {
			PrintEvent(event)
		}
	}
}

type EventBuilder struct {
	Event models.Event
}

func NewEventBuilder() *EventBuilder {
	return &EventBuilder{}
}

func (e *EventBuilder) AddName(name string) *EventBuilder {
	e.Event.Name = name
	return e
}

func (e *EventBuilder) AddDescription() *EventBuilder {
	var description string
	fmt.Println("enter description")
	fmt.Scanf("%s", description)
	e.Event.Description = description
	return e
}

func CreateNewEvent() models.Event {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter event name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter event description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Enter Number of artists : ")
	var numArtists int
	artists := make([]string, numArtists)
	fmt.Scanf("%d", &numArtists)
	for i := 0; i < numArtists; i++ {
		artist, _ := reader.ReadString('\n')
		artist = strings.TrimSpace(artist)
		artists[i] = artist
	}

	fmt.Print("Enter event duration (e.g. 2h30m): ")
	duration, _ := reader.ReadString('\n')
	duration = strings.TrimSpace(duration)
	var category models.EventCategory

	for {
		fmt.Println("Select category:")
		fmt.Printf("1. %s 2.%s 3.%s 4.%s 5.%s", models.Movie, models.Sports, models.Concert, models.Party, models.Workshop)
		catChoice, _ := reader.ReadString('\n')
		catChoice = strings.TrimSpace(catChoice)
		cat, _ := strconv.Atoi(strings.TrimSpace(catChoice))
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
