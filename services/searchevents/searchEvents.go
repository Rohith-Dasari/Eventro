package searchevents

import (
	"bufio"
	"eventro/config"
	"eventro/models"
	"eventro/storage"
	"fmt"
	"os"
	"strings"
)

func Search() {
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
		SearchByEventName()
	case 2:
		SearchByCategory()
	case 3:
		SearchByLocation()
	case 4:
		return
	}
}

func SearchByEventName() {

	fmt.Println("Enter the name of the event:")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	query := strings.TrimSpace(input)

	events := storage.LoadEvents()

	var found bool
	fmt.Println("Matching Events:")
	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Name), strings.ToLower(query)) {
			printEvent(event)
			found = true
		}
	}
	if !found {
		fmt.Println("No matching events found.")
	}
}

func SearchByCategory() {
	fmt.Println("Available Categories: movie, sports, concert, workshop, party")
	fmt.Print("Enter category: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	query := strings.ToLower(strings.TrimSpace(input))

	events := storage.LoadEvents()
	var found bool

	fmt.Println("Events in selected category:")
	for _, event := range events {
		if strings.ToLower(string(event.Category)) == query {
			printEvent(event)
			found = true
		}
	}
	if !found {
		fmt.Println("No events found in this category.")
	}
}

func SearchByLocation() {
	fmt.Print("Enter city to search for events: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	city := strings.ToLower(strings.TrimSpace(input))
	events := storage.LoadEvents()

	found := false
	fmt.Printf("\nEvents in %s:\n", city)
	for _, event := range events {
		if event.IsBlocked {
			continue
		}
		if contains(event.Locations, city) {
			printEvent(event)
			found = true
		}
	}

	if !found {
		fmt.Println("No events found in this city.")
	}

	// venues := storage.LoadVenues()
	// foundVenues := make([]string, len(venues)) //make map instead storing venue id and venue name
	// i := 0
	// //print venues
	// for _, venue := range venues {
	// 	if strings.ToLower(string(venue.City)) == query {
	// 		foundVenues[i] = venue.Name
	// 	}
	// 	i++
	// }
	// //print venue name and shows associated with it

	// //shows events whose shows are located in particular city
	// //first venues in city, find shows who has that venue and display them.
}
func contains(cities []string, target string) bool {
	for _, city := range cities {
		if strings.EqualFold(city, target) {
			return true
		}
	}
	return false

}

func printEvent(e models.Event) {
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
