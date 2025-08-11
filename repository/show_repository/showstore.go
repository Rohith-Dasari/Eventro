package showrepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
	"fmt"
	"log"
	"os"
)

type ShowRepository struct {
	Shows []models.Show
}

// func LoadShows() []models.Show {
// 	//read json
// 	data, err := os.ReadFile(config.ShowsFile)
// 	if err != nil {
// 		log.Fatalf("failed to read file %v", err)
// 	}
// 	//unmarshal into shows class
// 	var shows []models.Show
// 	if err := json.Unmarshal(data, &shows); err != nil {
// 		log.Fatalf("failed to marshal: %v", err)
// 	}
// 	return shows
// }

func (*ShowRepository) SaveShows(shows []models.Show) error {
	data, err := json.MarshalIndent(shows, "", "")
	if err != nil {
		return err
	}
	err = os.WriteFile(config.ShowsFile, data, 0644)
	return err
}
func NewShowRepository() *ShowRepository {
	//read json
	data, err := os.ReadFile(config.ShowsFile)
	if err != nil {
		log.Fatalf("failed to read file %v", err)
	}
	//unmarshal into shows class
	var shows []models.Show
	if err := json.Unmarshal(data, &shows); err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	return &ShowRepository{shows}

}
func (sr *ShowRepository) GetShows() ([]models.Show, error) {
	data, err := os.ReadFile(config.ShowsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read shows file: %w", err)
	}

	var shows []models.Show
	if err := json.Unmarshal(data, &shows); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shows: %w", err)
	}

	return shows, nil
}
