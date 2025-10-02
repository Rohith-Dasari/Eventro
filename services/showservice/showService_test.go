package showservice

import (
	"context"
	"errors"
	"eventro2/middleware"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"eventro2/mocks"

	"go.uber.org/mock/gomock"
)

func TestNewShowService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowRepo := mocks.NewMockShowRepository(ctrl)
	mockVenueRepo := mocks.NewMockVenueRepository(ctrl)
	mockBookRepo := mocks.NewMockBookingRepository(ctrl)
	mockEventRepo := mocks.NewMockEventRepository(ctrl)

	tests := []struct {
		name string
		args struct {
			showRepo  showrepository.ShowRepository
			venueRepo venuerepository.VenueRepository
			bookRepo  bookingrepository.BookingRepository
			eventRepo eventsrepository.EventRepository
		}
		want ShowService
	}{
		{
			name: "constructs service with all repos",
			args: struct {
				showRepo  showrepository.ShowRepository
				venueRepo venuerepository.VenueRepository
				bookRepo  bookingrepository.BookingRepository
				eventRepo eventsrepository.EventRepository
			}{
				showRepo:  mockShowRepo,
				venueRepo: mockVenueRepo,
				bookRepo:  mockBookRepo,
				eventRepo: mockEventRepo,
			},
			want: ShowService{
				ShowRepo:  mockShowRepo,
				VenueRepo: mockVenueRepo,
				BookRepo:  mockBookRepo,
				EventRepo: mockEventRepo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShowService(tt.args.showRepo, tt.args.venueRepo, tt.args.bookRepo, tt.args.eventRepo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewShowService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShowService_UpdateShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowRepo := mocks.NewMockShowRepository(ctrl)
	type args struct {
		ctx    context.Context
		showID string
		userID string
		update models.UpdateShowData
	}

	// Shared values
	showID := "show-1"
	userID := "host-1"
	now := time.Now()
	updatedPrice := 250.0
	updatedTime := "18:00"
	updatedBlocked := true

	originalShow := &models.Show{
		ID:          showID,
		HostID:      userID,
		VenueID:     "venue-1",
		EventID:     "event-1",
		CreatedAt:   now,
		IsBlocked:   false,
		Price:       100.0,
		ShowDate:    now,
		ShowTime:    "12:00",
		BookedSeats: []string{},
	}

	tests := []struct {
		name       string
		setupMocks func()
		args       args
		want       models.ShowResponse
		wantErr    bool
	}{
		{
			name: "successful update",
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID(showID).Return(originalShow, nil)
				updatedShow := *originalShow
				updatedShow.Price = updatedPrice
				updatedShow.ShowTime = updatedTime
				updatedShow.IsBlocked = updatedBlocked
				mockShowRepo.EXPECT().Update(&updatedShow).Return(nil)
			},
			args: args{
				ctx:    context.Background(),
				showID: showID,
				userID: userID,
				update: models.UpdateShowData{
					Price:     &updatedPrice,
					ShowTime:  &updatedTime,
					IsBlocked: &updatedBlocked,
				},
			},
			want: models.ShowResponse{
				ID:          showID,
				HostID:      userID,
				VenueID:     "venue-1",
				EventID:     "event-1",
				CreatedAt:   now,
				IsBlocked:   updatedBlocked,
				Price:       updatedPrice,
				ShowDate:    now,
				ShowTime:    updatedTime,
				BookedSeats: []string{},
			},
			wantErr: false,
		},
		{
			name: "unauthorized update",
			setupMocks: func() {
				unauthorizedShow := *originalShow
				unauthorizedShow.HostID = "other-user"
				mockShowRepo.EXPECT().GetByID(showID).Return(&unauthorizedShow, nil)
			},
			args: args{
				ctx:    context.Background(),
				showID: showID,
				userID: userID,
				update: models.UpdateShowData{Price: &updatedPrice},
			},
			want:    models.ShowResponse{},
			wantErr: true,
		},
		{
			name: "GetByID error",
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID(showID).Return(nil, fmt.Errorf("not found"))
			},
			args: args{
				ctx:    context.Background(),
				showID: showID,
				userID: userID,
				update: models.UpdateShowData{},
			},
			want:    models.ShowResponse{},
			wantErr: true,
		},
		{
			name: "Update error",
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID(showID).Return(originalShow, nil)
				updatedShow := *originalShow
				updatedShow.Price = updatedPrice
				mockShowRepo.EXPECT().Update(&updatedShow).Return(fmt.Errorf("update failed"))
			},
			args: args{
				ctx:    context.Background(),
				showID: showID,
				userID: userID,
				update: models.UpdateShowData{Price: &updatedPrice},
			},
			want:    models.ShowResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			s := &ShowService{
				ShowRepo: mockShowRepo,
			}
			got, err := s.UpdateShow(tt.args.ctx, tt.args.showID, tt.args.userID, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateShow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateShow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShowService_BrowseShows(t *testing.T) {

	type args struct {
		ctx    context.Context
		filter models.ShowFilter
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowRepo := mocks.NewMockShowRepository(ctrl)

	now := time.Now()

	sampleShows := []models.Show{
		{
			ID:          "show-1",
			HostID:      "host-1",
			VenueID:     "venue-1",
			EventID:     "event-1",
			CreatedAt:   now,
			IsBlocked:   false,
			Price:       120.5,
			ShowDate:    now,
			ShowTime:    "18:00",
			BookedSeats: []string{"A1", "A2"},
		},
	}

	expectedResponse := []models.ShowResponse{
		{
			ID:          "show-1",
			HostID:      "host-1",
			VenueID:     "venue-1",
			EventID:     "event-1",
			CreatedAt:   now,
			IsBlocked:   false,
			Price:       120.5,
			ShowDate:    now,
			ShowTime:    "18:00",
			BookedSeats: []string{"A1", "A2"},
		},
	}

	tests := []struct {
		name       string
		setupMocks func()
		args       args
		want       []models.ShowResponse
		wantErr    bool
	}{
		{
			name: "successful browse",
			setupMocks: func() {
				mockShowRepo.EXPECT().Find(models.ShowFilter{}).Return(sampleShows, nil)
			},
			args: args{
				ctx:    context.Background(),
				filter: models.ShowFilter{},
			},
			want:    expectedResponse,
			wantErr: false,
		},
		{
			name: "error fetching shows",
			setupMocks: func() {
				mockShowRepo.EXPECT().Find(models.ShowFilter{}).Return(nil, fmt.Errorf("database error"))
			},
			args: args{
				ctx:    context.Background(),
				filter: models.ShowFilter{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			s := &ShowService{
				ShowRepo: mockShowRepo,
			}
			got, err := s.BrowseShows(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrowseShows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrowseShows() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestShowService_CreateShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowRepo := mocks.NewMockShowRepository(ctrl)

	type fields struct {
		ShowRepo  showrepository.ShowRepository
		VenueRepo venuerepository.VenueRepository
		BookRepo  bookingrepository.BookingRepository
		EventRepo eventsrepository.EventRepository
	}
	type args struct {
		ctx      context.Context
		eventID  string
		venueID  string
		hostID   string
		price    float64
		showDate time.Time
		showTime string
	}

	timestamp := time.Now()

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       models.ShowResponse
		wantErr    bool
		setupMocks func()
	}{
		{
			name: "successful show creation",
			fields: fields{
				ShowRepo:  mockShowRepo,
				VenueRepo: nil,
				BookRepo:  nil,
				EventRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				eventID:  "event-1",
				venueID:  "venue-1",
				hostID:   "host-1",
				price:    299.99,
				showDate: timestamp,
				showTime: "18:00",
			},
			want: models.ShowResponse{
				HostID:      "host-1",
				VenueID:     "venue-1",
				EventID:     "event-1",
				IsBlocked:   false,
				Price:       299.99,
				ShowDate:    timestamp,
				ShowTime:    "18:00",
				BookedSeats: []string{},
			},
			wantErr: false,
			setupMocks: func() {
				mockShowRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "failure to create show",
			fields: fields{
				ShowRepo:  mockShowRepo,
				VenueRepo: nil,
				BookRepo:  nil,
				EventRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				eventID:  "event-2",
				venueID:  "venue-2",
				hostID:   "host-2",
				price:    199.00,
				showDate: timestamp,
				showTime: "20:00",
			},
			want:    models.ShowResponse{},
			wantErr: true,
			setupMocks: func() {
				mockShowRepo.EXPECT().Create(gomock.Any()).Return(errors.New("DB error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}
			s := &ShowService{
				ShowRepo:  tt.fields.ShowRepo,
				VenueRepo: tt.fields.VenueRepo,
				BookRepo:  tt.fields.BookRepo,
				EventRepo: tt.fields.EventRepo,
			}
			got, err := s.CreateShow(tt.args.ctx, tt.args.eventID, tt.args.venueID, tt.args.hostID, tt.args.price, tt.args.showDate, tt.args.showTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateShow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Match non-random parts of response
			if !tt.wantErr {
				if got.HostID != tt.want.HostID ||
					got.VenueID != tt.want.VenueID ||
					got.EventID != tt.want.EventID ||
					got.Price != tt.want.Price ||
					got.ShowDate != tt.want.ShowDate ||
					got.ShowTime != tt.want.ShowTime ||
					!reflect.DeepEqual(got.BookedSeats, tt.want.BookedSeats) {
					t.Errorf("CreateShow() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestShowService_DeleteShow(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowRepo := mocks.NewMockShowRepository(ctrl)

	type fields struct {
		ShowRepo showrepository.ShowRepository
	}
	field := fields{mockShowRepo}
	type args struct {
		ctx    context.Context
		showID string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       error
		wantErr    bool
		setupMocks func()
	}{
		{
			name:    "Admin deletes show",
			fields:  field,
			args:    args{AddUserToContext(ctx, "user-1", "Admin"), "show-1"},
			want:    nil,
			wantErr: false,
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID("show-1").Return(&models.Show{HostID: "user-2"}, nil)
				mockShowRepo.EXPECT().Delete(gomock.Any()).Return(nil)
			},
		},
		{
			name:    "show not found",
			fields:  field,
			args:    args{AddUserToContext(ctx, "user-1", "Host"), "show-2"},
			want:    fmt.Errorf("show not found"),
			wantErr: true,
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID("show-2").Return(nil, errors.New("show not found"))
			},
		},
		{
			name:    "unauthorised",
			fields:  field,
			args:    args{AddUserToContext(ctx, "user-1", "Customer"), "show-1"},
			want:    errors.New("unauthorized to delete this show"),
			wantErr: true,
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID("show-1").Return(&models.Show{HostID: "user-2"}, nil)
			},
		},
		{
			name:    "db error",
			fields:  field,
			args:    args{AddUserToContext(ctx, "user-1", "Admin"), "show-1"},
			want:    fmt.Errorf("failed to delete show"),
			wantErr: true,
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID("show-1").Return(&models.Show{HostID: "user-2"}, nil)
				mockShowRepo.EXPECT().Delete("show-1").Return(errors.New("DB error"))
			},
		},
		{
			name:    "Host tries to delete other host's show",
			fields:  field,
			args:    args{AddUserToContext(ctx, "user-1", "Host"), "show-1"},
			want:    errors.New("unauthorized to delete this show"),
			wantErr: true,
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID("show-1").Return(&models.Show{HostID: "user-2"}, nil)
			},
		},

		{
			name:    "Host deletes successfully",
			fields:  field,
			args:    args{AddUserToContext(ctx, "user-1", "Host"), "show-1"},
			want:    nil,
			wantErr: false,
			setupMocks: func() {
				mockShowRepo.EXPECT().GetByID("show-1").Return(&models.Show{HostID: "user-1"}, nil)
				mockShowRepo.EXPECT().Delete("show-1").Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			s := &ShowService{
				ShowRepo: mockShowRepo,
			}
			err := s.DeleteShow(tt.args.ctx, tt.args.showID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteShow() error = %v, wantErr %v", err, tt.want)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.want.Error()) {
				t.Errorf("DeleteShow() error = %v, wantErrMsg %v", err, tt.want)
			}
		})
	}

}

func AddUserToContext(ctx context.Context, userId string, userRole string) context.Context {
	ctx = context.WithValue(ctx, middleware.ContextUserIDKey, userId)
	ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, userRole)
	return ctx
}
