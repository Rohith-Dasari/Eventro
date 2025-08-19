package artistservice

import (
	"context"
	"eventro2/config"
	"eventro2/models"
	artistrepository "eventro2/repository/artists_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type Artistservice struct {
	ArtistRepo artistrepository.ArtistRepository
}

func NewArtistService(artistRepo artistrepository.ArtistRepository) Artistservice {
	return Artistservice{
		ArtistRepo: artistRepo,
	}
}

func (as *Artistservice) CreateArtist(ctx context.Context) {
	fmt.Println("Enter artist details:")

	var name string
	for {
		fmt.Print("Name: ")
		name = strings.TrimSpace(utils.ReadLine())
		if name == "" {
			color.Red("Artist name cannot be empty.")
			continue
		}
		break
	}

	var bio string
	for {
		fmt.Print("Bio: ")
		bio = strings.TrimSpace(utils.ReadLine())
		if bio == "" {
			color.Red("Artist bio cannot be empty.")
			continue
		}
		if len(bio) < 10 {
			color.Red("Bio should be at least 10 characters long.")
			continue
		}
		break
	}

	artist := models.Artist{
		ID:   uuid.New().String(),
		Name: name,
		Bio:  bio,
	}

	if err := as.ArtistRepo.Create(&artist); err != nil {
		color.Red("Failed to create artist: %v", err)
		return
	}

	color.Green("Artist created successfully.")
	fmt.Printf("%-8s : %s\n", "ID", artist.ID)
	fmt.Printf("%-8s : %s\n", "Name", artist.Name)
	fmt.Printf("%-8s : %s\n", "Bio", artist.Bio)
}

func (s *Artistservice) BrowseArtist(ctx context.Context) {
	fmt.Println("Enter artist name to search:")
	input := strings.ToLower(utils.ReadLine())

	matched, err := s.ArtistRepo.SearchByName(input)
	if err != nil {
		color.Red("Failed to search artists: %v", err)
		return
	}

	if len(matched) == 0 {
		color.Red("No artist found with that name.")
		return
	}
	color.Cyan("Matched Artists:")
	for _, a := range matched {
		fmt.Printf("%-8s : %s\n", "ID", a.ID)
		fmt.Printf("%-8s : %s\n", "Name", a.Name)
	}

	fmt.Println("Enter Artist ID to view full profile:")
	artistID := strings.TrimSpace(utils.ReadLine())
	if artistID == "" {
		color.Red("Artist ID cannot be empty.")
		return
	}

	artist, err := s.ArtistRepo.GetByID(artistID)
	if err != nil || artist == nil {
		color.Red("Artist not found.")
		return
	}

	fmt.Println(config.Dash)
	fmt.Printf("%-8s : %s\n", "Name", artist.Name)
	fmt.Printf("%-8s : %s\n", "Bio", artist.Bio)
	fmt.Println(config.Dash)

	events, err := s.ArtistRepo.GetEventsByArtistID(artistID)
	if err != nil {
		color.Red("Failed to fetch events for artist: %v", err)
		return
	}

	if len(events) == 0 {
		color.Yellow("No events found for this artist.")
		return
	}

	color.Green("Events by this artist:")
	for _, event := range events {
		fmt.Println(config.Dash)
		fmt.Printf("%-8s : %s\n", "ID", event.ID)
		fmt.Printf("%-8s : %s\n", "Name", event.Name)
		fmt.Printf("%-8s : %s\n", "Category", event.Category)
		fmt.Println(config.Dash)
	}
}
