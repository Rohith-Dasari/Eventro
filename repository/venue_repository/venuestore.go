package venuerepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
	"fmt"
	"log"
	"os"
	"sync"
)

type VenueRepository struct {
	Venues []models.Venue
}

var (
	venueRepoInstance *VenueRepository
	venueRepoOnce     sync.Once
	venueRepoMutex    sync.Mutex
)

// func LoadVenues() []models.Venue {
// 	file, err := os.Open(config.VenuesFile)
// 	if err != nil {
// 		fmt.Println("Error opening venues file:", err)
// 		return nil
// 	}
// 	defer file.Close()

// 	var venues []models.Venue
// 	decoder := json.NewDecoder(file)
// 	if err := decoder.Decode(&venues); err != nil {
// 		fmt.Println("Error decoding venues:", err)
// 		return nil
// 	}
// 	return venues
// }

func (vr *VenueRepository) GetVenues() ([]models.Venue, error) {
	data, err := os.ReadFile(config.VenuesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read venues file: %w", err)
	}

	var venues []models.Venue
	if err := json.Unmarshal(data, &venues); err != nil {
		return nil, fmt.Errorf("failed to unmarshal venues: %w", err)
	}

	return venues, nil
}

func (*VenueRepository) SaveVenues(venues []models.Venue) error {
	data, err := json.MarshalIndent(venues, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(config.VenuesFile, data, 0644)
	return err
}

func NewVenueRepository() *VenueRepository {
	//read json
	data, err := os.ReadFile(config.VenuesFile)
	if err != nil {
		log.Fatalf("failed to read file %v", err)
	}

	//unmarshal into user class
	var venues []models.Venue
	if err := json.Unmarshal(data, &venues); err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	return &VenueRepository{venues}
}
