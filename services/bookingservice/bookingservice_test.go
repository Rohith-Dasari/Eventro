package bookingservice

import (
	"context"
	"errors"
	"eventro2/mocks"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestNewBookingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookRepo := mocks.NewMockBookingRepository(ctrl)
	mockShowRepo := mocks.NewMockShowRepository(ctrl)
	mockVenueRepo := mocks.NewMockVenueRepository(ctrl)
	mockEventRepo := mocks.NewMockEventRepository(ctrl)

	type args struct {
		bRepo bookingrepository.BookingRepository
		sRepo showrepository.ShowRepository
		vRepo venuerepository.VenueRepository
		eRepo eventsrepository.EventRepository
	}
	tests := []struct {
		name string
		args args
		want BookingService
	}{
		{
			name: "constructs booking service correctly",
			args: args{
				bRepo: mockBookRepo,
				sRepo: mockShowRepo,
				vRepo: mockVenueRepo,
				eRepo: mockEventRepo,
			},
			want: BookingService{
				BookingRepo: mockBookRepo,
				ShowRepo:    mockShowRepo,
				venueRepo:   mockVenueRepo,
				eventRepo:   mockEventRepo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBookingService(tt.args.bRepo, tt.args.sRepo, tt.args.vRepo, tt.args.eRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookingService_AddBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockShowRepo := mocks.NewMockShowRepository(ctrl)

	baseShow := &models.Show{
		ID:          "show-1",
		HostID:      "host-1",
		BookedSeats: []string{"A1"},
		Price:       100,
		IsBlocked:   false,
		ShowDate:    time.Now(),
		ShowTime:    "18:00",
		EventID:     "event-1",
		VenueID:     "venue-1",
		CreatedAt:   time.Now(),
	}

	// Copy show to avoid mutation affecting other tests
	copyShow := *baseShow
	copyShow.BookedSeats = append([]string{}, baseShow.BookedSeats...)

	// Setup mocks
	mockShowRepo.EXPECT().
		GetByID("show-1").
		Return(&copyShow, nil)

	mockShowRepo.EXPECT().
		Update(gomock.Any()).
		Return(nil).
		Times(1)

	mockBookingRepo.EXPECT().
		Create(gomock.Any()).
		Return(nil).
		Times(1)

	service := &BookingService{
		BookingRepo: mockBookingRepo,
		ShowRepo:    mockShowRepo,
	}

	// Requested seats are valid and not booked
	requestedSeats := []string{"B2", "C3"}

	bookingResponse, err := service.AddBooking(context.Background(), "user-1", "show-1", requestedSeats)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bookingResponse == nil {
		t.Fatal("expected booking response, got nil")
	}

	// Validate response data
	if bookingResponse.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got %v", bookingResponse.UserID)
	}

	if bookingResponse.ShowID != "show-1" {
		t.Errorf("expected ShowID 'show-1', got %v", bookingResponse.ShowID)
	}

	if len(bookingResponse.Seats) != len(requestedSeats) {
		t.Errorf("expected seats length %d, got %d", len(requestedSeats), len(bookingResponse.Seats))
	}
}

func TestBookingService_isValidTicket(t *testing.T) {
	tests := []struct {
		name          string
		userTicket    string
		bookedTickets []string
		want          bool
	}{
		{
			name:          "valid ticket, not booked",
			userTicket:    "A1",
			bookedTickets: []string{"A2", "B3"},
			want:          true,
		},
		{
			name:          "ticket already booked",
			userTicket:    "A2",
			bookedTickets: []string{"A2", "B3"},
			want:          false,
		},
		{
			name:          "invalid format - row out of range",
			userTicket:    "Z10",
			bookedTickets: []string{},
			want:          false,
		},
		{
			name:          "invalid format - column out of range",
			userTicket:    "B11",
			bookedTickets: []string{},
			want:          false,
		},
		{
			name:          "valid lowercase input",
			userTicket:    "c4",
			bookedTickets: []string{"A1"},
			want:          true,
		},
		{
			name:          "duplicate in different case",
			userTicket:    "b3",
			bookedTickets: []string{"B3"},
			want:          false,
		},
		{
			name:          "invalid format - empty input",
			userTicket:    "",
			bookedTickets: []string{},
			want:          false,
		},
		{
			name:          "invalid format - no number",
			userTicket:    "A",
			bookedTickets: []string{},
			want:          false,
		},
		{
			name:          "invalid format - no letter",
			userTicket:    "12",
			bookedTickets: []string{},
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BookingService{} // Mocks not needed here
			if got := bs.isValidTicket(tt.userTicket, tt.bookedTickets); got != tt.want {
				t.Errorf("isValidTicket(%q, %v) = %v, want %v", tt.userTicket, tt.bookedTickets, got, tt.want)
			}
		})
	}
}

func TestBookingService_AddBooking_BlockedShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockShowRepo := mocks.NewMockShowRepository(ctrl)

	blockedShow := &models.Show{
		ID:          "show-1",
		HostID:      "host-1",
		BookedSeats: []string{"A1"},
		Price:       100,
		IsBlocked:   true, // blocked show
		ShowDate:    time.Now(),
		ShowTime:    "18:00",
		EventID:     "event-1",
		VenueID:     "venue-1",
		CreatedAt:   time.Now(),
	}

	mockShowRepo.EXPECT().
		GetByID("show-1").
		Return(blockedShow, nil)

	service := &BookingService{
		BookingRepo: mockBookingRepo,
		ShowRepo:    mockShowRepo,
	}

	requestedSeats := []string{"B2"}

	_, err := service.AddBooking(context.Background(), "user-1", "show-1", requestedSeats)

	if err == nil {
		t.Fatal("expected error due to blocked show, got nil")
	}

	if err.Error() != "cannot book tickets for a blocked show" {
		t.Errorf("expected blocked show error, got %v", err)
	}
}

func TestBookingService_BrowseBookings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)

	type fields struct {
		BookingRepo bookingrepository.BookingRepository
		ShowRepo    showrepository.ShowRepository
		venueRepo   venuerepository.VenueRepository
		eventRepo   eventsrepository.EventRepository
	}
	type args struct {
		ctx       context.Context
		bookingID string
		userID    string
		showID    string
	}
	now := time.Now()

	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      []models.BookingResponse
		wantErr   bool
	}{
		{
			name: "successfully returns bookings for user",
			fields: fields{
				BookingRepo: mockBookingRepo,
			},
			args: args{
				ctx:    context.Background(),
				userID: "user-123",
			},
			mockSetup: func() {
				mockBookingRepo.EXPECT().
					Find("", "user-123", "").
					Return([]models.Booking{
						{
							BookingID:         "b1",
							UserID:            "user-123",
							ShowID:            "show-1",
							TimeBooked:        now,
							NumTickets:        2,
							TotalBookingPrice: 300,
							Seats:             []string{"A1", "A2"},
						},
					}, nil)
			},
			want: []models.BookingResponse{
				{
					BookingID:         "b1",
					UserID:            "user-123",
					ShowID:            "show-1",
					TimeBooked:        now,
					NumTickets:        2,
					TotalBookingPrice: 300,
					Seats:             []string{"A1", "A2"},
				},
			},
			wantErr: false,
		},
		{
			name: "returns error when repository fails",
			fields: fields{
				BookingRepo: mockBookingRepo,
			},
			args: args{
				ctx:    context.Background(),
				userID: "user-error",
			},
			mockSetup: func() {
				mockBookingRepo.EXPECT().
					Find("", "user-error", "").
					Return(nil, errors.New("database failure"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			s := &BookingService{
				BookingRepo: tt.fields.BookingRepo,
				ShowRepo:    tt.fields.ShowRepo,
				venueRepo:   tt.fields.venueRepo,
				eventRepo:   tt.fields.eventRepo,
			}
			got, err := s.BrowseBookings(tt.args.ctx, tt.args.bookingID, tt.args.userID, tt.args.showID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrowseBookings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) || (len(got) > 0 && got[0].BookingID != tt.want[0].BookingID) {
				t.Errorf("BrowseBookings() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
