package userrepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
	"fmt"
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

func (ur *UserRepository) UserExists(users []models.User, email string) bool {
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (ur *UserRepository) GetUsers() ([]models.User, error) {
	data, err := os.ReadFile(config.UsersFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read users file: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	return users, nil
}

func (*UserRepository) SaveUsers(users []models.User) error {
	//how to marshal slice to json?
	//use ident
	data, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialise data: %w", err)
	}
	err = os.WriteFile(config.UsersFile, data, 0644)
	return err
}

func (ur *UserRepository) AddUser(user models.User) error {
	ur.Users = append(ur.Users, user)
	return ur.SaveUsers(ur.Users)
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
