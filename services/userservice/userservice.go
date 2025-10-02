package userservice

import (
	"context"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"fmt"
)

type UserService struct {
	UserRepo userrepository.UserRepository
}

func NewUserService(userRepo userrepository.UserRepository) UserService {
	return UserService{UserRepo: userRepo}
}

func (s *UserService) BrowseUsers(ctx context.Context, userID string, blocked *bool) ([]models.User, error) {
	if blocked != nil && *blocked {
		return s.UserRepo.GetBlockedUsers()
	}

	return s.UserRepo.GetUsers()
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	if userID != "" {
		return s.UserRepo.GetByID(userID)
	}
	return nil, fmt.Errorf("invalid user ID")
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, req models.UpdateUserRequest) (models.User, error) {
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return models.User{}, err
	}

	if req.IsBlocked != nil {
		user.IsBlocked = *req.IsBlocked
	}

	if req.Role != nil {
		user.Role = models.Role(*req.Role)
	}

	if err := s.UserRepo.Update(user); err != nil {
		return models.User{}, err
	}

	return *user, nil
}
