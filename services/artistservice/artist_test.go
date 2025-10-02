package artistservice_test

import (
	"context"
	"errors"
	"eventro2/models"
	"eventro2/services/artistservice"
	"strings"
	"testing"
)

type mockArtistRepo struct {
	ReturnCreateError           bool
	ReturnDeleteError           bool
	ReturnUpdateError           bool
	ReturnGetByIDError          bool
	ReturnListError             bool
	ReturnListByEventError      bool
	ReturnSearchByNameError     bool
	ReturnEventsByArtistIDError bool

	CreatedArtist            *models.Artist
	UpdatedArtist            *models.Artist
	DeletedArtistID          string
	GetByIDInput             string
	ListByEventInput         string
	GetEventsByArtistIDInput string
	SearchByNameInput        string
}

func (m *mockArtistRepo) Create(artist *models.Artist) error {
	m.CreatedArtist = artist
	if m.ReturnCreateError {
		return errors.New("create error")
	}
	return nil
}

func (m *mockArtistRepo) GetByID(id string) (*models.Artist, error) {
	m.GetByIDInput = id
	if m.ReturnGetByIDError {
		return nil, errors.New("get by id error")
	}
	return &models.Artist{ID: id, Name: "Mock Artist", Bio: "Mock bio"}, nil
}

func (m *mockArtistRepo) List() ([]models.Artist, error) {
	if m.ReturnListError {
		return nil, errors.New("list error")
	}
	return []models.Artist{
		{ID: "1", Name: "A", Bio: "B"},
	}, nil
}

func (m *mockArtistRepo) ListByEvent(eventID string) ([]models.Artist, error) {
	m.ListByEventInput = eventID
	if m.ReturnListByEventError {
		return nil, errors.New("list by event error")
	}
	return []models.Artist{
		{ID: "2", Name: "Event Artist", Bio: "Performer bio"},
	}, nil
}

func (m *mockArtistRepo) Update(artist *models.Artist) error {
	m.UpdatedArtist = artist
	if m.ReturnUpdateError {
		return errors.New("update error")
	}
	return nil
}

func (m *mockArtistRepo) Delete(id string) error {
	m.DeletedArtistID = id
	if m.ReturnDeleteError {
		return errors.New("delete error")
	}
	return nil
}

func (m *mockArtistRepo) GetEventsByArtistID(artistID string) ([]models.Event, error) {
	m.GetEventsByArtistIDInput = artistID
	if m.ReturnEventsByArtistIDError {
		return nil, errors.New("get events error")
	}
	return []models.Event{
		{ID: "e1", Name: "Event 1", Description: "A great event"},
	}, nil
}

func (m *mockArtistRepo) SearchByName(name string) ([]models.Artist, error) {
	m.SearchByNameInput = name
	if m.ReturnSearchByNameError {
		return nil, errors.New("search error")
	}
	return []models.Artist{
		{ID: "3", Name: name, Bio: "Bio for " + name},
	}, nil
}
func TestCreateArtist(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		inputBio      string
		mockReturnErr bool
		wantErr       bool
		errMsg        string
	}{
		{
			name:      "success case",
			inputName: "rohith dasari",
			inputBio:  "This is a valid biography",
		},
		{
			name:      "bio too short",
			inputName: "sai ganesh",
			inputBio:  "short bio",
			wantErr:   true,
			errMsg:    "bio must be at least 12 characters long",
		},
		{
			name:          "repository create error",
			inputName:     "Error Artist",
			inputBio:      "This biography is definitely long enough",
			mockReturnErr: true,
			wantErr:       true,
			errMsg:        "create error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArtistRepo{
				ReturnCreateError: tt.mockReturnErr,
			}
			svc := artistservice.NewArtistService(mockRepo)

			artist, err := svc.CreateArtist(context.Background(), tt.inputName, tt.inputBio)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if artist.Name != tt.inputName {
				t.Errorf("expected artist name %q, got %q", tt.inputName, artist.Name)
			}

			if artist.Bio != tt.inputBio {
				t.Errorf("expected artist bio %q, got %q", tt.inputBio, artist.Bio)
			}

			if mockRepo.CreatedArtist == nil {
				t.Errorf("expected Create to be called on repository")
			}
		})
	}
}

func TestDeleteArtist(t *testing.T) {
	tests := []struct {
		name          string
		artistID      string
		mockReturnErr bool
		wantErr       bool
		errMsg        string
	}{
		{
			name:     "successful delete",
			artistID: "artist-123",
		},
		{
			name:          "repository delete error",
			artistID:      "artist-error",
			mockReturnErr: true,
			wantErr:       true,
			errMsg:        "delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArtistRepo{
				ReturnDeleteError: tt.mockReturnErr,
			}
			svc := artistservice.NewArtistService(mockRepo)

			err := svc.DeleteArtist(context.Background(), tt.artistID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if mockRepo.DeletedArtistID != tt.artistID {
				t.Errorf("expected artist ID %q to be deleted, got %q", tt.artistID, mockRepo.DeletedArtistID)
			}
		})
	}
}

func TestGetArtists(t *testing.T) {
	tests := []struct {
		name                  string
		searchName            string
		returnSearchByNameErr bool
		returnListErr         bool
		wantErr               bool
		errMsg                string
	}{
		{
			name:       "search by name success",
			searchName: "John",
		},
		{
			name:                  "search by name error",
			searchName:            "Jane",
			returnSearchByNameErr: true,
			wantErr:               true,
			errMsg:                "search error",
		},
		{
			name:          "list all artists success",
			searchName:    "",
			returnListErr: false,
		},
		{
			name:          "list all artists error",
			searchName:    "",
			returnListErr: true,
			wantErr:       true,
			errMsg:        "list error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArtistRepo{
				ReturnSearchByNameError: tt.returnSearchByNameErr,
				ReturnListError:         tt.returnListErr,
			}
			svc := artistservice.NewArtistService(mockRepo)

			_, err := svc.GetArtists(context.Background(), tt.searchName)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.searchName != "" && mockRepo.SearchByNameInput != tt.searchName {
				t.Errorf("expected SearchByName to be called with %q, got %q", tt.searchName, mockRepo.SearchByNameInput)
			}
		})
	}
}

func TestArtistService_GetArtistByID(t *testing.T) {
	tests := []struct {
		name             string
		inputID          string
		returnGetByIDErr bool
		wantErr          bool
		errMsg           string
	}{
		{
			name:    "successfully gets artist by ID",
			inputID: "artist-123",
		},
		{
			name:             "repository returns error",
			inputID:          "non-existent-id",
			returnGetByIDErr: true,
			wantErr:          true,
			errMsg:           "get by id error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArtistRepo{
				ReturnGetByIDError: tt.returnGetByIDErr,
			}
			service := artistservice.NewArtistService(mockRepo)

			artist, err := service.GetArtistByID(context.Background(), tt.inputID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if artist.ID != tt.inputID {
				t.Errorf("expected artist ID %q, got %q", tt.inputID, artist.ID)
			}

			if mockRepo.GetByIDInput != tt.inputID {
				t.Errorf("expected repo GetByID called with %q, got %q", tt.inputID, mockRepo.GetByIDInput)
			}
		})
	}
}
