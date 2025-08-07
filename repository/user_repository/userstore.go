package userrepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
	"log"
	"os"
)

type UserRepository struct {
	Users []models.User
}

// func LoadUsers() []models.User {
// 	//read json
// 	data, err := os.ReadFile("data/users.json")
// 	if err != nil {
// 		log.Fatalf("failed to read file %v", err)
// 	}

// 	//unmarshal into user class
// 	var users []models.User
// 	if err := json.Unmarshal(data, &users); err != nil {
// 		log.Fatalf("failed to marshal: %v", err)
// 	}
// 	return users
// }

// func UserExists(users []models.User, email string) bool {
// 	for _, user := range users {
// 		if user.Email == email {
// 			return true
// 		}
// 	}
// 	return false
// }

func (*UserRepository) SaveUsers(users []models.User) error {
	//how to marshal slice to json?
	//use ident
	data, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(config.VenuesFile, data, 0644)
	return err
}
func NewUserRepository() *UserRepository {
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

	return &UserRepository{users}
}
