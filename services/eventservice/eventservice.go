package eventservice

import (
	"context"
	"eventro/models"
	"eventro/storage"
	"fmt"
	"strings"
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
