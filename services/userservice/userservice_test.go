package userservice

import (
	"context"
	"eventro2/mocks"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	tests := []struct {
		name string
		args userrepository.UserRepository
		want UserService
	}{
		{
			name: "construct with mock user repo",
			args: mockUserRepo,
			want: UserService{
				UserRepo: mockUserRepo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserService(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_BrowseUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	x := true

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserRepo.EXPECT().GetUsers().Return([]models.User{}, nil)
	mockUserRepo.EXPECT().GetBlockedUsers().Return([]models.User{}, nil)
	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		ctx     context.Context
		userID  string
		blocked *bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.User
		wantErr bool
	}{
		{
			name: "browse all users",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:     context.Background(),
				userID:  "user-1",
				blocked: nil,
			},
			want:    []models.User{},
			wantErr: false,
		},
		{
			name: "blocked users",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:     context.Background(),
				userID:  "user-1",
				blocked: &x,
			},
			want:    []models.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{
				UserRepo: tt.fields.UserRepo,
			}
			got, err := s.BrowseUsers(tt.args.ctx, tt.args.userID, tt.args.blocked)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.BrowseUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.BrowseUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserRepo.EXPECT().GetByID("user-1").Return(&models.User{}, nil).MaxTimes(1)
	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "valid user ID", fields: fields{
			UserRepo: mockUserRepo,
		}, args: args{ctx: context.Background(), userID: "user-1"}, want: &models.User{}, wantErr: false},
		{name: "invalid user ID", fields: fields{
			UserRepo: mockUserRepo,
		}, args: args{ctx: context.Background(), userID: ""}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{
				UserRepo: tt.fields.UserRepo,
			}
			got, err := s.GetUserByID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := &UserService{
		UserRepo: mockRepo,
	}

	userID := "user-1"
	existingUser := &models.User{
		UserID:    userID,
		IsBlocked: false,
		Role:      "user",
	}

	t.Run("successful update of IsBlocked", func(t *testing.T) {
		newBlocked := true
		req := models.UpdateUserRequest{
			IsBlocked: &newBlocked,
		}

		updatedUser := *existingUser
		updatedUser.IsBlocked = newBlocked

		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().Update(&updatedUser).Return(nil)

		got, err := service.UpdateUser(context.Background(), userID, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got.IsBlocked != newBlocked {
			t.Errorf("expected IsBlocked = %v, got %v", newBlocked, got.IsBlocked)
		}
	})

	t.Run("successful update of Role", func(t *testing.T) {
		newRole := "admin"
		req := models.UpdateUserRequest{
			Role: &newRole,
		}

		updatedUser := *existingUser
		updatedUser.Role = models.Role(newRole)

		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().Update(&updatedUser).Return(nil)

		got, err := service.UpdateUser(context.Background(), userID, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got.Role != models.Role(newRole) {
			t.Errorf("expected Role = %v, got %v", newRole, got.Role)
		}
	})

	t.Run("update both fields", func(t *testing.T) {
		newBlocked := false
		newRole := "moderator"
		req := models.UpdateUserRequest{
			IsBlocked: &newBlocked,
			Role:      &newRole,
		}

		updatedUser := *existingUser
		updatedUser.IsBlocked = newBlocked
		updatedUser.Role = models.Role(newRole)

		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().Update(&updatedUser).Return(nil)

		got, err := service.UpdateUser(context.Background(), userID, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got.IsBlocked != newBlocked || got.Role != models.Role(newRole) {
			t.Errorf("update failed: got %+v", got)
		}
	})

	t.Run("GetByID returns error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(userID).Return(nil, fmt.Errorf("not found"))

		_, err := service.UpdateUser(context.Background(), userID, models.UpdateUserRequest{})
		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})

	t.Run("Update returns error", func(t *testing.T) {
		newBlocked := true
		req := models.UpdateUserRequest{
			IsBlocked: &newBlocked,
		}

		updatedUser := *existingUser
		updatedUser.IsBlocked = newBlocked

		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().Update(&updatedUser).Return(fmt.Errorf("update failed"))

		_, err := service.UpdateUser(context.Background(), userID, req)
		if err == nil {
			t.Errorf("expected update error but got nil")
		}
	})
}
