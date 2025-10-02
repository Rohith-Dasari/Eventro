package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"eventro2/handlers"
	"eventro2/middleware"
	"eventro2/mocks"
	"eventro2/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestBookingHandler_CreateBooking(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		role           string
		userID         string
		requestBody    handlers.CreateBookingRequest
		mockSetup      func(m *mocks.MockBookingServiceInterface)
		expectedStatus int
	}{
		{
			name:   "success - customer creates booking",
			method: http.MethodPost,
			role:   "Customer",
			userID: "user-123",
			requestBody: handlers.CreateBookingRequest{
				ShowID: "show-456",
				Seats:  []string{"A1", "A2"},
			},
			mockSetup: func(m *mocks.MockBookingServiceInterface) {
				m.EXPECT().
					AddBooking(gomock.Any(), "user-123", "show-456", []string{"A1", "A2"}).
					Return(&models.BookingResponse{BookingID: "b1"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "invalid method",
			method: http.MethodGet,
			role:   "Customer",
			userID: "user-123",
			requestBody: handlers.CreateBookingRequest{
				ShowID: "show-456",
				Seats:  []string{"A1"},
			},
			mockSetup:      func(m *mocks.MockBookingServiceInterface) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "forbidden role",
			method: http.MethodPost,
			role:   "Guest",
			userID: "user-123",
			requestBody: handlers.CreateBookingRequest{
				ShowID: "show-456",
				Seats:  []string{"A1"},
			},
			mockSetup:      func(m *mocks.MockBookingServiceInterface) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid request body",
			method:         http.MethodPost,
			role:           "Customer",
			userID:         "user-123",
			requestBody:    handlers.CreateBookingRequest{},
			mockSetup:      func(m *mocks.MockBookingServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookingServiceInterface(ctrl)
			tt.mockSetup(mockService)

			handler := &handlers.BookingHandler{
				BookingService: mockService,
			}

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(tt.method, "/bookings", bytes.NewReader(bodyBytes))

			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.CreateBooking(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
func TestBookingHandler_BrowseBookings(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		role           string
		userID         string
		query          string
		mockSetup      func(m *mocks.MockBookingServiceInterface)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:   "success - admin browsing all",
			method: http.MethodGet,
			role:   "admin",
			userID: "admin-1",
			query:  "?showId=show-123",
			mockSetup: func(m *mocks.MockBookingServiceInterface) {
				m.EXPECT().
					BrowseBookings(gomock.Any(), "", "", "show-123").
					Return([]models.BookingResponse{{BookingID: "b1"}, {BookingID: "b2"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:   "success - customer only their bookings",
			method: http.MethodGet,
			role:   "Customer",
			userID: "user-123",
			query:  "",
			mockSetup: func(m *mocks.MockBookingServiceInterface) {
				m.EXPECT().
					BrowseBookings(gomock.Any(), "", "user-123", "").
					Return([]models.BookingResponse{{BookingID: "b3"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			role:           "Customer",
			userID:         "user-123",
			query:          "",
			mockSetup:      func(m *mocks.MockBookingServiceInterface) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCount:  0,
		},
		{
			name:   "service error",
			method: http.MethodGet,
			role:   "Customer",
			userID: "user-123",
			query:  "?bookingId=b1",
			mockSetup: func(m *mocks.MockBookingServiceInterface) {
				m.EXPECT().
					BrowseBookings(gomock.Any(), "b1", "user-123", "").
					Return(nil, assertAnError())
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookingServiceInterface(ctrl)
			tt.mockSetup(mockService)

			handler := &handlers.BookingHandler{
				BookingService: mockService,
			}

			req := httptest.NewRequest(tt.method, "/bookings"+tt.query, nil)

			ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tt.userID)
			ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.BrowseBookings(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if rr.Code == http.StatusOK {
				var bookings []models.BookingResponse
				if err := json.NewDecoder(rr.Body).Decode(&bookings); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(bookings) != tt.expectedCount {
					t.Errorf("expected %d bookings, got %d", tt.expectedCount, len(bookings))
				}
			}
		})
	}
}

func assertAnError() error {
	return fmt.Errorf("some error")
}
