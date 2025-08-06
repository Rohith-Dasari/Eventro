package storage

import (
	"encoding/json"
	"eventro/config"
	"eventro/models"
	"fmt"
	"os"
)

func LoadVenues() []models.Venue {
	file, err := os.Open(config.VenuesFile)
	if err != nil {
		fmt.Println("Error opening venues file:", err)
		return nil
	}
	defer file.Close()

	var venues []models.Venue
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&venues); err != nil {
		fmt.Println("Error decoding venues:", err)
		return nil
	}
	return venues
}
