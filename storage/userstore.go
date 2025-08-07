package storage

import (
	"encoding/json"
	"eventro/config"
	"eventro/models"
	"log"
	"os"
)

func LoadUsers() []models.User {
	//read json
	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Fatalf("failed to read file %v", err)
	}

	//unmarshal into user class
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	return users
}

func UserExists(users []models.User, email string) bool {
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func SaveUsers(users []models.User) error {
	//how to marshal slice to json?
	//use ident
	data, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(config.VenuesFile, data, 0644)
	return err

}
