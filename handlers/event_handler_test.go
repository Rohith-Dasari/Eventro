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

func TestEventHandler_CreateEvent(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		input          CreateEventRequest
		mockSetup      func(m *mocks.MockEventServiceI)
		expectedStatus int
	}{
		{
			name: "success - admin creates event",
			role: "Admin",
			input: CreateEventRequest{
				Name:        "Concert",
				Description: "Music Show",
				Duration:    "2h",
				Category:    "Music",
				Artists:     []string{"artist1"},
			},
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().
					CreateNewEvent(gomock.Any(), "Concert", "Music Show", "2h", models.EventCategory("Music"), []string{"artist1"}).
					Return(models.EventResponse{ID: "event1", Name: "Concert"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "forbidden - non-admin user",
			role:           "Customer",
			input:          CreateEventRequest{Name: "Concert"},
			mockSetup:      func(m *mocks.MockEventServiceI) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:  "invalid request body",
			role:  "Admin",
			input: CreateEventRequest{},
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().
					CreateNewEvent(gomock.Any(), "", "", "", models.EventCategory(""), []string(nil)).
					Return(models.EventResponse{}, errors.New("invalid input"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockEventServiceI(ctrl)
			tt.mockSetup(mockService)

			handler := &EventHandler{EventService: mockService}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.CreateEvent(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
func TestEventHandler_BrowseEvents(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		query          string
		mockSetup      func(m *mocks.MockEventServiceI)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "invalid method",
			method: http.MethodPost,
			query:  "",
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "success - returns events",
			method: http.MethodGet,
			query:  "?eventname=Concert",
			mockSetup: func(m *mocks.MockEventServiceI) {
				expectedFilter := models.EventFilter{
					EventID:    "",
					Name:       "Concert",
					Category:   "",
					Location:   "",
					IsBlocked:  nil,
					ArtistName: "",
				}
				m.EXPECT().
					BrowseEvents(gomock.Any(), expectedFilter).
					Return([]models.EventResponse{{ID: "1", Name: "Concert"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []models.EventResponse{{ID: "1", Name: "Concert"}},
		},
		{
			name:   "service error",
			method: http.MethodGet,
			query:  "?category=Music",
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().
					BrowseEvents(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventService := mocks.NewMockEventServiceI(ctrl)
			tt.mockSetup(mockEventService)

			h := &EventHandler{EventService: mockEventService}

			req := httptest.NewRequest(tt.method, "/events"+tt.query, nil)
			w := httptest.NewRecorder()

			h.BrowseEvents(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedBody != nil && res.StatusCode == http.StatusOK {
				var got []models.EventResponse
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.expectedBody)
				if !bytes.Equal(gotJSON, wantJSON) {
					t.Errorf("expected body %s, got %s", string(wantJSON), string(gotJSON))
				}
			}
		})
	}
}

func TestEventHandler_DeleteEvent(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		eventID        string
		role           string
		mockSetup      func(m *mocks.MockEventServiceI)
		expectedStatus int
	}{
		{
			name:    "invalid method",
			method:  http.MethodGet,
			eventID: "123",
			role:    "Admin",
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:    "forbidden - non-admin user",
			method:  http.MethodDelete,
			eventID: "123",
			role:    "User",
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:    "invalid request - missing eventID",
			method:  http.MethodDelete,
			eventID: "",
			role:    "Admin",
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "internal server error - service failure",
			method:  http.MethodDelete,
			eventID: "123",
			role:    "Admin",
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().DeleteEvent(gomock.Any(), "123").Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "success - event deleted",
			method:  http.MethodDelete,
			eventID: "123",
			role:    "Admin",
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().DeleteEvent(gomock.Any(), "123").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventService := mocks.NewMockEventServiceI(ctrl)
			tt.mockSetup(mockEventService)

			h := &EventHandler{EventService: mockEventService}

			req := httptest.NewRequest(tt.method, "/events/"+tt.eventID, nil)
			req.SetPathValue("eventID", tt.eventID)

			ctx := context.WithValue(req.Context(), middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeleteEvent(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestEventHandler_UpdateEvent(t *testing.T) {
	name := "updated"
	tests := []struct {
		name           string
		method         string
		eventID        string
		body           interface{}
		mockSetup      func(m *mocks.MockEventServiceI)
		expectedStatus int
	}{
		{
			name:    "invalid method",
			method:  http.MethodGet,
			eventID: "123",
			body:    models.EventUpdate{},
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:    "missing eventID",
			method:  http.MethodPatch,
			eventID: "",
			body:    models.EventUpdate{},
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "invalid json body",
			method:  http.MethodPatch,
			eventID: "123",
			body:    "not-json",
			mockSetup: func(m *mocks.MockEventServiceI) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "service error",
			method:  http.MethodPatch,
			eventID: "123",
			body:    models.EventUpdate{Name: &name},
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().UpdateEvent(gomock.Any(), "123", gomock.Any()).Return(models.EventResponse{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "success",
			method:  http.MethodPatch,
			eventID: "123",
			body:    models.EventUpdate{Name: &name},
			mockSetup: func(m *mocks.MockEventServiceI) {
				m.EXPECT().UpdateEvent(gomock.Any(), "123", gomock.Any()).Return(models.EventResponse{ID: "123", Name: "Updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventService := mocks.NewMockEventServiceI(ctrl)
			tt.mockSetup(mockEventService)

			h := &EventHandler{EventService: mockEventService}

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

			req := httptest.NewRequest(tt.method, "/events/"+tt.eventID, bytes.NewReader(bodyBytes))
			req.SetPathValue("eventID", tt.eventID)

			w := httptest.NewRecorder()
			h.UpdateEvent(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}
