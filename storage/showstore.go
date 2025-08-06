package storage

import (
	"encoding/json"
	"eventro/config"
	"eventro/models"
	"log"
	"os"
)

func LoadShows() []models.Show {
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
	return shows
}
