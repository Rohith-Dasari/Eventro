package eventservice

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	eventsrepository "eventro2/repository/event_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type EventService struct {
	EventRepo eventsrepository.EventRepository
}

func NewEventService(eventRepo eventsrepository.EventRepository) EventService {
	return EventService{EventRepo: eventRepo}
}

func (e *EventService) ModerateEvents(ctx context.Context) {
	for {
		fmt.Println("Enter ID of the event to be moderated:")
		eventID := utils.ReadLine()

		requiredEvent, err := e.EventRepo.GetByID(eventID)
		if err != nil {
			fmt.Println("Error fetching event:", err)
			return
		}
		if requiredEvent == nil {
			fmt.Println("Event not found. Please enter a correct ID.")
			fmt.Println("1. Retry with another Event ID\n2. Back ")
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

		// Print current event details
		e.PrintEvent(*requiredEvent)

		if requiredEvent.IsBlocked {
			fmt.Println("The event is currently BLOCKED")
			fmt.Println("Are you sure you want to UNBLOCK the event?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")

			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				requiredEvent.IsBlocked = false
				if err := e.EventRepo.Update(requiredEvent); err != nil {
					color.Red("Failed to unblock event: %v", err)
				} else {
					color.Green("Event is successfully unblocked")
				}
				return
			case 2:
				color.Red("Canceled")
				continue
			case 3:
				return
			}
		} else {
			fmt.Println("The event is currently UNBLOCKED")
			fmt.Println("Are you sure you want to BLOCK the event?")
			fmt.Println("1. Yes")
			fmt.Println("2. No, try another ID")
			fmt.Println("3. Back")

			choice, _ := utils.TakeUserInput()
			switch choice {
			case 1:
				requiredEvent.IsBlocked = true
				if err := e.EventRepo.Update(requiredEvent); err != nil {
					color.Red("Failed to block event: %v", err)
				} else {
					fmt.Println("Event is successfully blocked")
				}
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

func (e *EventService) ViewBlockedEvents(ctx context.Context) {
	events, err := e.EventRepo.List()
	if err != nil {
		color.Red("Error fetching events: %v", err)
		return
	}

	found := false
	for _, event := range events {
		if event.IsBlocked {
			e.PrintEvent(event)
			found = true
		}
	}
	if !found {
		color.Red("No Blocked Events")
	}
}
func (e *EventService) CreateNewEvent() {
	reader := bufio.NewReader(os.Stdin)

	// Name
	var name string
	for {
		fmt.Print("Enter event name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			color.Red("Event name cannot be empty.")
			continue
		}
		break
	}

	// Description
	var description string
	for {
		fmt.Print("Enter event description: ")
		description, _ = reader.ReadString('\n')
		description = strings.TrimSpace(description)
		if description == "" {
			color.Red("Description cannot be empty.")
			continue
		}
		if len(description) < 10 {
			color.Red("Description should be at least 10 characters long.")
			continue
		}
		break
	}

	// Number of artists
	var numArtists int
	for {
		var err error
		fmt.Print("Enter number of artists: ")
		numArtists, err = utils.TakeUserInput()
		if err != nil {
			color.Red("Please enter a valid integer.")
			continue
		}
		if numArtists <= 0 {
			color.Red("There must be at least one artist.")
			continue
		}
		break
	}

	// Artists
	artists := make([]string, numArtists)
	for i := 0; i < numArtists; i++ {
		for {
			fmt.Printf("Enter artist %d: ", i+1)
			ID := utils.ReadLine()
			if ID == "" {
				color.Red("Artist ID cannot be empty.")
				continue
			}
			artists[i] = ID
			break
		}
	}

	// Duration
	var duration string
	for {
		fmt.Print("Enter event duration (e.g., 2h30m): ")
		duration, _ = reader.ReadString('\n')
		duration = strings.TrimSpace(duration)
		if duration == "" {
			color.Red("Duration cannot be empty.")
			continue
		}
		_, err := time.ParseDuration(duration)
		if err != nil {
			color.Red("Invalid duration format. Use formats like 1h, 2h30m, 45m.")
			continue
		}
		break
	}

	// Category
	var category models.EventCategory
	for {
		fmt.Println("Select category:")
		fmt.Printf("1. %s\n2. %s\n3. %s\n4. %s\n5. %s\n",
			models.Movie, models.Sports, models.Concert, models.Party, models.Workshop)

		catChoice, _ := reader.ReadString('\n')
		catChoice = strings.TrimSpace(catChoice)
		cat, err := strconv.Atoi(catChoice)

		if err != nil || cat < 1 || cat > 5 {
			color.Red("Invalid choice. Please enter a number between 1 and 5.")
			continue
		}

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
		}
		break
	}

	EventID := uuid.New().String()

	//populate eventartist table with eventID and artist id stored in array
	for i := range artists {
		eventartist := models.EventArtist{EventID: EventID, ArtistID: artists[i]}
		e.EventRepo.AddEventArtist(&eventartist)
	}

	// Build event
	event := models.Event{
		ID:          EventID,
		Name:        name,
		Description: description,
		Duration:    duration,
		Category:    category,
		IsBlocked:   false,
	}

	if err := e.EventRepo.Create(&event); err != nil {
		log.Printf("Error saving event: %v", err)
		return
	}

	e.PrintEvent(event)
	color.Green("Event added successfully")
}

func (e *EventService) PrintEvent(event models.Event) {
	fmt.Println(config.Dash)
	fmt.Printf("ID: %s\n", event.ID)
	fmt.Printf("Name: %s\n", event.Name)
	fmt.Printf("Description: %s\n", event.Description)
	fmt.Printf("Hype Meter: %d\n", event.HypeMeter)
	fmt.Printf("Duration: %s\n", event.Duration)
	fmt.Printf("Category: %s\n", event.Category)

	artists, err := e.EventRepo.GetArtistsByEventID(event.ID)
	if err != nil {
		fmt.Println("Error fetching artists:", err)
		return
	}
	fmt.Println("Artists:")
	for _, artist := range artists {
		fmt.Printf(" - %s\n", artist.Name)
	}

	fmt.Println(config.Dash)
}
