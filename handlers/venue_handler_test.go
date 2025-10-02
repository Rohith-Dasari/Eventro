package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"eventro2/middleware"
	"eventro2/mocks"
	"eventro2/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVenueHandler_BrowseVenues(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		query          string
		mockSetup      func(m *mocks.MockVenueServiceInterface)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			query:          "",
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "success - returns venues",
			method: http.MethodGet,
			query:  "?city=NYC",
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				filter := models.VenueFilter{
					City:      "NYC",
					HostID:    "",
					VenueID:   "",
					IsBlocked: false,
				}
				m.EXPECT().
					BrowseVenues(gomock.Any(), filter).
					Return([]models.VenueResponse{
						{ID: "v1", Name: "Madison Square Garden", City: "NYC", State: "NY", IsSeatLayoutRequired: true},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []models.VenueResponse{
				{ID: "v1", Name: "Madison Square Garden", City: "NYC", State: "NY", IsSeatLayoutRequired: true},
			},
		},
		{
			name:   "service error",
			method: http.MethodGet,
			query:  "?city=LA",
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				filter := models.VenueFilter{
					City:      "LA",
					HostID:    "",
					VenueID:   "",
					IsBlocked: false,
				}
				m.EXPECT().
					BrowseVenues(gomock.Any(), filter).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockVenueServiceInterface(ctrl)
			tt.mockSetup(mockService)

			h := &VenueHandler{VenueService: mockService}

			req := httptest.NewRequest(tt.method, "/venues"+tt.query, nil)
			w := httptest.NewRecorder()

			h.BrowseVenues(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != nil && res.StatusCode == http.StatusOK {
				var got []models.VenueResponse
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}
}

func TestVenueHandler_DeleteVenue(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		venueID        string
		userID         string
		userRole       string
		mockSetup      func(m *mocks.MockVenueServiceInterface)
		expectedStatus int
	}{
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			venueID:        "v1",
			userID:         "u1",
			userRole:       "Host",
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "missing venueID",
			method:         http.MethodDelete,
			venueID:        "",
			userID:         "u1",
			userRole:       "Host",
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized user",
			method:         http.MethodDelete,
			venueID:        "v1",
			userID:         "u1",
			userRole:       "Customer",
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:     "service error",
			method:   http.MethodDelete,
			venueID:  "v1",
			userID:   "u1",
			userRole: "Host",
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				m.EXPECT().
					DeleteVenue(gomock.Any(), "v1", "u1", "host").
					Return(assert.AnError)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:     "success",
			method:   http.MethodDelete,
			venueID:  "v1",
			userID:   "u1",
			userRole: "Admin",
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				m.EXPECT().
					DeleteVenue(gomock.Any(), "v1", "u1", "admin").
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockVenueServiceInterface(ctrl)
			tt.mockSetup(mockService)

			h := &VenueHandler{VenueService: mockService}

			req := httptest.NewRequest(tt.method, "/venues/"+tt.venueID, nil)
			if tt.venueID != "" {
				req.SetPathValue("venueID", tt.venueID)
			}

			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.userRole)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeleteVenue(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestVenueHandler_CreateVenue(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		userRole       string
		body           interface{}
		mockSetup      func(m *mocks.MockVenueServiceInterface)
		expectedStatus int
		expectedBody   *models.VenueResponse
	}{
		{
			name:           "invalid method",
			method:         http.MethodGet,
			userID:         "host1",
			userRole:       "Host",
			body:           CreateVenueRequest{},
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "unauthorized role",
			method:         http.MethodPost,
			userID:         "host1",
			userRole:       "Customer",
			body:           CreateVenueRequest{},
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing userID",
			method:         http.MethodPost,
			userID:         "",
			userRole:       "Host",
			body:           CreateVenueRequest{},
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid JSON body",
			method:         http.MethodPost,
			userID:         "host1",
			userRole:       "Host",
			body:           "not-json",
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			method:   http.MethodPost,
			userID:   "host1",
			userRole: "Host",
			body: CreateVenueRequest{
				Name:                 "Venue1",
				City:                 "CityA",
				State:                "StateA",
				IsSeatLayoutRequired: true,
			},
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				m.EXPECT().
					CreateVenue(gomock.Any(), "host1", "Venue1", "CityA", "StateA", true).
					Return(models.VenueResponse{}, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "success",
			method:   http.MethodPost,
			userID:   "host1",
			userRole: "Host",
			body: CreateVenueRequest{
				Name:                 "Venue1",
				City:                 "CityA",
				State:                "StateA",
				IsSeatLayoutRequired: true,
			},
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				m.EXPECT().
					CreateVenue(gomock.Any(), "host1", "Venue1", "CityA", "StateA", true).
					Return(models.VenueResponse{
						ID:                   "v1",
						Name:                 "Venue1",
						City:                 "CityA",
						State:                "StateA",
						IsSeatLayoutRequired: true,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &models.VenueResponse{
				ID:                   "v1",
				Name:                 "Venue1",
				City:                 "CityA",
				State:                "StateA",
				IsSeatLayoutRequired: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockVenueServiceInterface(ctrl)
			tt.mockSetup(mockService)

			h := &VenueHandler{VenueService: mockService}

			var bodyBytes []byte
			var err error
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/venues", bytes.NewReader(bodyBytes))
			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.userRole)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreateVenue(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != nil && res.StatusCode == http.StatusCreated {
				var got models.VenueResponse
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				assert.Equal(t, *tt.expectedBody, got)
			}
		})
	}
}
func TestVenueHandler_UpdateVenue(t *testing.T) {
	name := "Updated Venue"
	city := "New City"
	isBlocked := true
	isSeatLayout := true

	tests := []struct {
		name           string
		method         string
		venueID        string
		userID         string
		userRole       string
		body           interface{}
		mockSetup      func(m *mocks.MockVenueServiceInterface)
		expectedStatus int
		expectedBody   *models.VenueResponse
	}{
		{
			name:           "invalid method",
			method:         http.MethodGet,
			venueID:        "v1",
			userID:         "u1",
			userRole:       "Host",
			body:           UpdateVenueRequest{},
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "missing venueID",
			method:         http.MethodPatch,
			venueID:        "",
			userID:         "u1",
			userRole:       "Host",
			body:           UpdateVenueRequest{},
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized user",
			method:         http.MethodPatch,
			venueID:        "v1",
			userID:         "u1",
			userRole:       "Customer",
			body:           UpdateVenueRequest{},
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid JSON body",
			method:         http.MethodPatch,
			venueID:        "v1",
			userID:         "u1",
			userRole:       "Host",
			body:           "not-json",
			mockSetup:      func(m *mocks.MockVenueServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			method:   http.MethodPatch,
			venueID:  "v1",
			userID:   "u1",
			userRole: "Host",
			body: UpdateVenueRequest{
				Name: &name,
			},
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				m.EXPECT().
					UpdateVenue(gomock.Any(), "v1", "u1", "Host", gomock.Any()).
					Return(models.VenueResponse{}, assert.AnError)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:     "success",
			method:   http.MethodPatch,
			venueID:  "v1",
			userID:   "u1",
			userRole: "Host",
			body: UpdateVenueRequest{
				Name:                 &name,
				City:                 &city,
				IsBlocked:            &isBlocked,
				IsSeatLayoutRequired: &isSeatLayout,
			},
			mockSetup: func(m *mocks.MockVenueServiceInterface) {
				m.EXPECT().
					UpdateVenue(gomock.Any(), "v1", "u1", "Host", gomock.Any()).
					Return(models.VenueResponse{
						ID:                   "v1",
						Name:                 name,
						City:                 city,
						IsBlocked:            isBlocked,
						IsSeatLayoutRequired: isSeatLayout,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &models.VenueResponse{
				ID:                   "v1",
				Name:                 name,
				City:                 city,
				IsBlocked:            isBlocked,
				IsSeatLayoutRequired: isSeatLayout,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockVenueServiceInterface(ctrl)
			tt.mockSetup(mockService)

			h := &VenueHandler{VenueService: mockService}

			var bodyBytes []byte
			var err error
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/venues/"+tt.venueID, bytes.NewReader(bodyBytes))
			if tt.venueID != "" {
				req.SetPathValue("venueID", tt.venueID)
			}

			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.userRole)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdateVenue(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != nil && res.StatusCode == http.StatusOK {
				var got models.VenueResponse
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				assert.Equal(t, *tt.expectedBody, got)
			}
		})
	}
}
