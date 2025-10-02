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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func ptrFloat64(f float64) *float64 { return &f }
func ptrString(s string) *string    { return &s }
func ptrBool(b bool) *bool          { return &b }

func TestShowHandler_BrowseShows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShowServiceInterface(ctrl)
	h := &ShowHandler{ShowService: mockService}

	tests := []struct {
		name           string
		method         string
		query          string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "successful browse",
			method: http.MethodGet,
			query:  "?eventId=e1",
			mockSetup: func() {
				mockService.EXPECT().
					BrowseShows(gomock.Any(), models.ShowFilter{EventID: "e1"}).
					Return([]models.ShowResponse{{ID: "s1"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			query:          "?eventId=e1",
			mockSetup:      func() {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "service error",
			method: http.MethodGet,
			query:  "?eventId=e2",
			mockSetup: func() {
				mockService.EXPECT().
					BrowseShows(gomock.Any(), models.ShowFilter{EventID: "e2"}).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			req := httptest.NewRequest(tt.method, "/shows"+tt.query, nil)
			w := httptest.NewRecorder()

			h.BrowseShows(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestShowHandler_UpdateShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowService := mocks.NewMockShowServiceInterface(ctrl)
	handler := &ShowHandler{ShowService: mockShowService}

	tests := []struct {
		name           string
		showID         string
		userID         string
		userRole       string
		requestBody    UpdateShowRequest
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:     "successful update",
			showID:   "s1",
			userID:   "host1",
			userRole: "Host",
			requestBody: UpdateShowRequest{
				Price:     ptrFloat64(250),
				ShowTime:  ptrString("18:00"),
				IsBlocked: ptrBool(true),
			},
			mockSetup: func() {
				expectedResponse := models.ShowResponse{
					ID:        "s1",
					HostID:    "host1",
					VenueID:   "v1",
					EventID:   "e1",
					Price:     250,
					ShowTime:  "18:00",
					IsBlocked: true,
				}
				mockShowService.EXPECT().
					UpdateShow(gomock.Any(), "s1", "host1", gomock.Any()).
					Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "unauthorized user",
			showID:   "s1",
			userID:   "host2",
			userRole: "Host",
			requestBody: UpdateShowRequest{
				Price: ptrFloat64(250),
			},
			mockSetup: func() {
				mockShowService.EXPECT().
					UpdateShow(gomock.Any(), "s1", "host2", gomock.Any()).
					Return(models.ShowResponse{}, errors.New("forbidden"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "invalid method",
			showID:         "s1",
			userID:         "host1",
			userRole:       "Host",
			requestBody:    UpdateShowRequest{},
			mockSetup:      func() {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:     "invalid show date format",
			showID:   "s2",
			userID:   "host1",
			userRole: "Host",
			requestBody: UpdateShowRequest{
				ShowDate: ptrString("2025-13-99"), // invalid date
			},
			mockSetup:      func() {}, // service should not be called
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			method := http.MethodPatch
			if tt.expectedStatus == http.StatusMethodNotAllowed {
				method = http.MethodGet
			}

			req := httptest.NewRequest(method, "/shows/"+tt.showID, bytes.NewReader(bodyBytes))
			req.SetPathValue("showID", tt.showID)
			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.userRole)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.UpdateShow(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
func TestShowHandler_CreateShow(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		input          CreateShowRequest
		mockSetup      func(m *mocks.MockShowServiceInterface)
		expectedStatus int
	}{
		{
			name: "success - host creates show",
			role: "Host",
			input: CreateShowRequest{
				EventID:  "event1",
				VenueID:  "venue1",
				Price:    200,
				ShowDate: "2025-08-28",
				ShowTime: "18:00",
			},
			mockSetup: func(m *mocks.MockShowServiceInterface) {
				parsedDate, _ := time.Parse("2006-01-02", "2025-08-28")
				m.EXPECT().
					CreateShow(gomock.Any(), "event1", "venue1", "host1", 200.0, parsedDate, "18:00").
					Return(models.ShowResponse{
						ID:          "show1",
						HostID:      "host1",
						EventID:     "event1",
						VenueID:     "venue1",
						Price:       200,
						ShowDate:    parsedDate,
						ShowTime:    "18:00",
						IsBlocked:   false,
						BookedSeats: []string{},
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "forbidden - non-host user",
			role:           "Customer",
			input:          CreateShowRequest{EventID: "event1"},
			mockSetup:      func(m *mocks.MockShowServiceInterface) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid request body",
			role: "Host",
			input: CreateShowRequest{
				EventID: "",
			},
			mockSetup:      func(m *mocks.MockShowServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid date format",
			role: "Host",
			input: CreateShowRequest{
				EventID:  "event1",
				VenueID:  "venue1",
				ShowDate: "28-08-2025",
				ShowTime: "18:00",
			},
			mockSetup:      func(m *mocks.MockShowServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error - CreateShow fails",
			role: "Host",
			input: CreateShowRequest{
				EventID:  "event1",
				VenueID:  "venue1",
				Price:    200,
				ShowDate: "2025-08-28",
				ShowTime: "18:00",
			},
			mockSetup: func(m *mocks.MockShowServiceInterface) {
				parsedDate, _ := time.Parse("2006-01-02", "2025-08-28")
				m.EXPECT().
					CreateShow(gomock.Any(), "event1", "venue1", "host1", 200.0, parsedDate, "18:00").
					Return(models.ShowResponse{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockShowServiceInterface(ctrl)
			tt.mockSetup(mockService)

			handler := &ShowHandler{ShowService: mockService}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/shows", bytes.NewBuffer(body))

			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, "host1")
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.CreateShow(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestShowHandler_DeleteShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowService := mocks.NewMockShowServiceInterface(ctrl)
	handler := &ShowHandler{ShowService: mockShowService}

	tests := []struct {
		name           string
		showID         string
		userID         string
		userRole       string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:     "missing showID",
			showID:   "",
			userID:   "host1",
			userRole: "Host",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			showID:   "s1",
			userID:   "host1",
			userRole: "Host",
			mockSetup: func() {
				mockShowService.EXPECT().DeleteShow(gomock.Any(), "s1").Return(errors.New("DB error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "successful delete",
			showID:   "s1",
			userID:   "host1",
			userRole: "Host",
			mockSetup: func() {
				mockShowService.EXPECT().DeleteShow(gomock.Any(), "s1").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodDelete, "/shows/"+tt.showID, nil)
			req.SetPathValue("showID", tt.showID)

			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.userRole)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.DeleteShow(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestShowHandler_DeleteShow_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowService := mocks.NewMockShowServiceInterface(ctrl)
	handler := &ShowHandler{ShowService: mockShowService}

	req := httptest.NewRequest(http.MethodGet, "/shows/s1", nil)
	rr := httptest.NewRecorder()

	handler.DeleteShow(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestShowHandler_BrowseShows_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShowService := mocks.NewMockShowServiceInterface(ctrl)
	handler := &ShowHandler{ShowService: mockShowService}

	// Setup: service returns an error
	mockShowService.EXPECT().
		BrowseShows(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/shows", nil)
	rr := httptest.NewRecorder()

	handler.BrowseShows(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	expected := "Failed to fetch shows: database error"
	body := rr.Body.String()
	if !strings.Contains(body, expected) {
		t.Errorf("expected body to contain %q, got %q", expected, body)
	}
}
