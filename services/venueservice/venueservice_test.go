package venueservice

import (
	"context"
	"eventro2/mocks"
	"eventro2/models"
	venuerepository "eventro2/repository/venue_repository"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestNewVenueService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		repo venuerepository.VenueRepository
	}
	tests := []struct {
		name string
		args args
		want VenueService
	}{
		// TODO: Add test cases.
		{
			name: "valid repo",
			args: args{
				repo: mocks.NewMockVenueRepository(ctrl),
			},
			want: VenueService{
				VenueRepo: mocks.NewMockVenueRepository(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVenueService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVenueService() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestVenueService_CreateVenue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		VenueRepo venuerepository.VenueRepository
	}
	type args struct {
		ctx                  context.Context
		hostID               string
		name                 string
		city                 string
		state                string
		isSeatLayoutRequired bool
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(mockRepo *mocks.MockVenueRepository)
		wantErr   bool
	}{
		{
			name: "successful creation",
			fields: fields{
				VenueRepo: mocks.NewMockVenueRepository(ctrl),
			},
			args: args{
				ctx:                  context.Background(),
				hostID:               "host-1",
				name:                 "Venue 1",
				city:                 "City 1",
				state:                "State 1",
				isSeatLayoutRequired: true,
			},
			mockSetup: func(mockRepo *mocks.MockVenueRepository) {
				mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "unsuccessful creation",
			fields: fields{
				VenueRepo: mocks.NewMockVenueRepository(ctrl),
			},
			args: args{
				ctx:                  context.Background(),
				hostID:               "host-1",
				name:                 "Venue 1",
				city:                 "City 1",
				state:                "State 1",
				isSeatLayoutRequired: true,
			},
			mockSetup: func(mockRepo *mocks.MockVenueRepository) {
				mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.VenueRepo.(*mocks.MockVenueRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			vs := &VenueService{
				VenueRepo: mockRepo,
			}

			got, err := vs.CreateVenue(tt.args.ctx, tt.args.hostID, tt.args.name, tt.args.city, tt.args.state, tt.args.isSeatLayoutRequired)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVenue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.HostID != tt.args.hostID || got.Name != tt.args.name || got.City != tt.args.city || got.State != tt.args.state || got.IsSeatLayoutRequired != tt.args.isSeatLayoutRequired {
					t.Errorf("CreateVenue() returned incorrect data: %+v", got)
				}
			}
		})
	}
}
func TestVenueService_UpdateVenue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockVenueRepository(ctrl)
	service := &VenueService{
		VenueRepo: mockRepo,
	}

	existingVenue := &models.Venue{
		ID:                   "venue-1",
		HostID:               "host-1",
		Name:                 "Old Name",
		City:                 "Old City",
		State:                "Old State",
		IsSeatLayoutRequired: false,
		IsBlocked:            false,
	}

	type args struct {
		ctx      context.Context
		venueID  string
		userID   string
		userRole string
		update   models.UpdateVenueData
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func()
		want      models.VenueResponse
		wantErr   bool
	}{
		{
			name: "host successfully updates venue name",
			args: args{
				ctx:      context.Background(),
				venueID:  "venue-1",
				userID:   "host-1",
				userRole: "user",
				update: models.UpdateVenueData{
					Name: strPtr("New Name"),
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(cloneVenue(existingVenue), nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			want: models.VenueResponse{
				ID:                   "venue-1",
				HostID:               "host-1",
				Name:                 "New Name",
				City:                 "Old City",
				State:                "Old State",
				IsSeatLayoutRequired: false,
			},
			wantErr: false,
		},
		{
			name: "non-host cannot update venue",
			args: args{
				ctx:      context.Background(),
				venueID:  "venue-1",
				userID:   "user-2",
				userRole: "user",
				update: models.UpdateVenueData{
					Name: strPtr("Invalid Update"),
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(cloneVenue(existingVenue), nil)
			},
			want:    models.VenueResponse{},
			wantErr: true,
		},
		{
			name: "admin successfully blocks venue",
			args: args{
				ctx:      context.Background(),
				venueID:  "venue-1",
				userID:   "admin-1",
				userRole: "admin",
				update: models.UpdateVenueData{
					IsBlocked: boolPtr(true),
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(cloneVenue(existingVenue), nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			want: models.VenueResponse{
				ID:                   "venue-1",
				HostID:               "host-1",
				Name:                 "Old Name",
				City:                 "Old City",
				State:                "Old State",
				IsSeatLayoutRequired: false,
			},
			wantErr: false,
		},
		{
			name: "non-admin tries to block venue",
			args: args{
				ctx:      context.Background(),
				venueID:  "venue-1",
				userID:   "random-user",
				userRole: "user",
				update: models.UpdateVenueData{
					IsBlocked: boolPtr(true),
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(cloneVenue(existingVenue), nil)
			},
			want:    models.VenueResponse{},
			wantErr: true,
		},
		{
			name: "update fails due to repo error",
			args: args{
				ctx:      context.Background(),
				venueID:  "venue-1",
				userID:   "host-1",
				userRole: "user",
				update: models.UpdateVenueData{
					State: strPtr("New State"),
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(cloneVenue(existingVenue), nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(fmt.Errorf("update error"))
			},
			want:    models.VenueResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.UpdateVenue(tt.args.ctx, tt.args.venueID, tt.args.userID, tt.args.userRole, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVenue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateVenue() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// Helper functions
func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

// Clone to avoid mutation across tests
func cloneVenue(v *models.Venue) *models.Venue {
	copy := *v
	return &copy
}
func TestVenueService_DeleteVenue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockVenueRepository(ctrl)
	service := &VenueService{VenueRepo: mockRepo}

	existingVenue := &models.Venue{
		ID:     "venue-1",
		HostID: "host-1",
		Name:   "Test Venue",
	}

	tests := []struct {
		name      string
		venueID   string
		userID    string
		userRole  string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:     "host successfully deletes venue",
			venueID:  "venue-1",
			userID:   "host-1",
			userRole: "user",
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(existingVenue, nil)
				mockRepo.EXPECT().Delete("venue-1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "admin successfully deletes venue",
			venueID:  "venue-1",
			userID:   "admin-1",
			userRole: "admin",
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(existingVenue, nil)
				mockRepo.EXPECT().Delete("venue-1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "non-host and non-admin cannot delete venue",
			venueID:  "venue-1",
			userID:   "user-2",
			userRole: "user",
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(existingVenue, nil)
			},
			wantErr: true,
		},
		{
			name:     "GetByID returns error",
			venueID:  "venue-1",
			userID:   "host-1",
			userRole: "user",
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(nil, fmt.Errorf("venue not found"))
			},
			wantErr: true,
		},
		{
			name:     "Delete returns error",
			venueID:  "venue-1",
			userID:   "host-1",
			userRole: "user",
			mockSetup: func() {
				mockRepo.EXPECT().GetByID("venue-1").Return(existingVenue, nil)
				mockRepo.EXPECT().Delete("venue-1").Return(fmt.Errorf("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.DeleteVenue(context.Background(), tt.venueID, tt.userID, tt.userRole)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteVenue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVenueService_BrowseVenues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockVenueRepository(ctrl)
	service := &VenueService{
		VenueRepo: mockRepo,
	}

	sampleVenues := []models.Venue{
		{
			ID:                   "venue-1",
			HostID:               "host-1",
			Name:                 "Test Venue 1",
			City:                 "City A",
			State:                "State X",
			IsSeatLayoutRequired: true,
		},
		{
			ID:                   "venue-2",
			HostID:               "host-2",
			Name:                 "Test Venue 2",
			City:                 "City B",
			State:                "State Y",
			IsSeatLayoutRequired: false,
		},
	}

	expectedResponse := []models.VenueResponse{
		{
			ID:                   "venue-1",
			HostID:               "host-1",
			Name:                 "Test Venue 1",
			City:                 "City A",
			State:                "State X",
			IsSeatLayoutRequired: true,
		},
		{
			ID:                   "venue-2",
			HostID:               "host-2",
			Name:                 "Test Venue 2",
			City:                 "City B",
			State:                "State Y",
			IsSeatLayoutRequired: false,
		},
	}

	tests := []struct {
		name      string
		filter    models.VenueFilter
		mockSetup func()
		want      []models.VenueResponse
		wantErr   bool
	}{
		{
			name:   "successful fetch",
			filter: models.VenueFilter{}, // could be empty or have values
			mockSetup: func() {
				mockRepo.EXPECT().Find(models.VenueFilter{}).Return(sampleVenues, nil)
			},
			want:    expectedResponse,
			wantErr: false,
		},
		{
			name:   "repository returns error",
			filter: models.VenueFilter{City: "Unknown"},
			mockSetup: func() {
				mockRepo.EXPECT().Find(models.VenueFilter{City: "Unknown"}).Return(nil, fmt.Errorf("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.BrowseVenues(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrowseVenues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrowseVenues() = %v, want %v", got, tt.want)
			}
		})
	}
}
