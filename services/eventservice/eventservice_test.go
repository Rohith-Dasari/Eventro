package eventservice

import (
	"context"
	"errors"
	"eventro2/mocks"
	"eventro2/models"
	eventsrepository "eventro2/repository/event_repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestNewEventService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepository(ctrl)

	type args struct {
		eventRepo eventsrepository.EventRepository
	}
	tests := []struct {
		name string
		args args
		want EventService
	}{
		{
			name: "valid constructor",
			args: args{
				eventRepo: mockEventRepo,
			},
			want: EventService{
				EventRepo: mockEventRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventService(tt.args.eventRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventService_CreateNewEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)

	type fields struct {
		EventRepo eventsrepository.EventRepository
	}
	type args struct {
		ctx         context.Context
		name        string
		description string
		duration    string
		category    models.EventCategory
		artistIDs   []string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      models.EventResponse
		wantErr   bool
	}{
		{
			name: "successful event creation",
			fields: fields{
				EventRepo: mockRepo,
			},
			args: args{
				ctx:         context.Background(),
				name:        "Rock Night",
				description: "A night of rock music",
				duration:    "2h",
				category:    models.Concert,
				artistIDs:   []string{"artist1", "artist2"},
			},
			mockSetup: func() {
				// Simulate successful creation
				mockRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

				// Simulate AddEventArtist for each artist
				mockRepo.EXPECT().AddEventArtist(gomock.Any()).Return(nil).Times(2)
			},
			want: models.EventResponse{
				Name:        "Rock Night",
				Description: "A night of rock music",
				Duration:    "2h",
				Category:    string(models.Concert),
				IsBlocked:   false,
				ArtistIDs:   []string{"artist1", "artist2"},
			},
			wantErr: false,
		},
		{
			name: "event creation failure",
			fields: fields{
				EventRepo: mockRepo,
			},
			args: args{
				ctx:         context.Background(),
				name:        "Test Event",
				description: "Testing failure",
				duration:    "1h",
				category:    models.Movie,
				artistIDs:   []string{},
			},
			mockSetup: func() {
				mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error")).Times(1)
			},
			want:    models.EventResponse{},
			wantErr: true,
		},
		{
			name: "artist association fails",
			fields: fields{
				EventRepo: mockRepo,
			},
			args: args{
				ctx:         context.Background(),
				name:        "Test Event 2",
				description: "Some event",
				duration:    "90m",
				category:    models.Movie,
				artistIDs:   []string{"a1"},
			},
			mockSetup: func() {
				mockRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
				mockRepo.EXPECT().AddEventArtist(gomock.Any()).Return(errors.New("artist association fail")).Times(1)
			},
			want:    models.EventResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			service := &EventService{
				EventRepo: tt.fields.EventRepo,
			}
			got, err := service.CreateNewEvent(tt.args.ctx, tt.args.name, tt.args.description, tt.args.duration, tt.args.category, tt.args.artistIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNewEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// only compare dynamic ID if no error
			if !tt.wantErr {
				tt.want.ID = got.ID // accept generated ID
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNewEvent() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
func TestEventService_BrowseEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepository(ctrl)

	type fields struct {
		EventRepo eventsrepository.EventRepository
	}
	type args struct {
		ctx    context.Context
		filter models.EventFilter
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      []models.EventResponse
		wantErr   bool
	}{
		{
			name: "successfully returns filtered events",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx: context.Background(),
				filter: models.EventFilter{
					Category: string(models.Concert),
				},
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					GetFilteredEvents(models.EventFilter{Category: string(models.Concert)}).
					Return([]models.Event{
						{
							ID:          "e1",
							Name:        "Rock Show",
							Description: "Live Rock",
							Duration:    "2h",
							Category:    models.Concert,
							IsBlocked:   false,
						},
						{
							ID:          "e2",
							Name:        "Jazz Evening",
							Description: "Smooth jazz",
							Duration:    "90m",
							Category:    models.Concert,
							IsBlocked:   false,
						},
					}, nil).
					Times(1)
			},
			want: []models.EventResponse{
				{
					ID:          "e1",
					Name:        "Rock Show",
					Description: "Live Rock",
					Duration:    "2h",
					Category:    "concert",
					IsBlocked:   false,
				},
				{
					ID:          "e2",
					Name:        "Jazz Evening",
					Description: "Smooth jazz",
					Duration:    "90m",
					Category:    "concert",
					IsBlocked:   false,
				},
			},
			wantErr: false,
		},
		{
			name: "repository returns error",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx:    context.Background(),
				filter: models.EventFilter{},
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					GetFilteredEvents(gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
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

			s := &EventService{
				EventRepo: tt.fields.EventRepo,
			}
			got, err := s.BrowseEvents(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrowseEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrowseEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestEventService_DeleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepository(ctrl)

	type fields struct {
		EventRepo eventsrepository.EventRepository
	}
	type args struct {
		ctx     context.Context
		eventID string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful deletion",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx:     context.Background(),
				eventID: "event-123",
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					Delete("event-123").
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "repository returns error on delete",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx:     context.Background(),
				eventID: "event-404",
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					Delete("event-404").
					Return(errors.New("delete failed")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			e := &EventService{
				EventRepo: tt.fields.EventRepo,
			}
			if err := e.DeleteEvent(tt.args.ctx, tt.args.eventID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepository(ctrl)

	type fields struct {
		EventRepo eventsrepository.EventRepository
	}
	type args struct {
		ctx        context.Context
		eventID    string
		updateData models.EventUpdate
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      models.EventResponse
		wantErr   bool
	}{
		{
			name: "successfully updates event name and duration",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx:     context.Background(),
				eventID: "event-123",
				updateData: models.EventUpdate{
					Name:     ptr("Updated Event"),
					Duration: ptr("3h"),
				},
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					GetByID("event-123").
					Return(&models.Event{
						ID:          "event-123",
						Name:        "Old Name",
						Description: "Desc",
						Duration:    "2h",
						Category:    models.Concert,
						IsBlocked:   false,
					}, nil)

				mockEventRepo.EXPECT().
					Update(gomock.AssignableToTypeOf(&models.Event{})).
					DoAndReturn(func(e *models.Event) error {
						if e.ID != "event-123" {
							return errors.New("wrong event ID")
						}
						if e.Name != "Updated Event" {
							return errors.New("name not updated")
						}
						if e.Duration != "3h" {
							return errors.New("duration not updated")
						}
						return nil
					})
			},
			want: models.EventResponse{
				ID:          "event-123",
				Name:        "Updated Event",
				Description: "Desc",
				Duration:    "3h",
				Category:    "concert",
				IsBlocked:   false,
			},
			wantErr: false,
		},
		{
			name: "event not found",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx:        context.Background(),
				eventID:    "missing-id",
				updateData: models.EventUpdate{Name: ptr("New")},
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					GetByID("missing-id").
					Return(nil, errors.New("not found"))
			},
			want:    models.EventResponse{},
			wantErr: true,
		},
		{
			name: "update fails",
			fields: fields{
				EventRepo: mockEventRepo,
			},
			args: args{
				ctx:     context.Background(),
				eventID: "event-123",
				updateData: models.EventUpdate{
					Description: ptr("new desc"),
				},
			},
			mockSetup: func() {
				mockEventRepo.EXPECT().
					GetByID("event-123").
					Return(&models.Event{
						ID:          "event-123",
						Name:        "Name",
						Description: "old desc",
						Duration:    "1h",
						Category:    models.Movie,
						IsBlocked:   false,
					}, nil)

				mockEventRepo.EXPECT().
					Update(gomock.AssignableToTypeOf(&models.Event{})).
					Return(errors.New("update error"))
			},
			want:    models.EventResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			s := &EventService{
				EventRepo: tt.fields.EventRepo,
			}
			got, err := s.UpdateEvent(tt.args.ctx, tt.args.eventID, tt.args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateEvent() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
