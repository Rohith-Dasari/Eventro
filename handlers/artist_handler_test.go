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

	gomock "go.uber.org/mock/gomock"
)

func TestArtistHandler_CreateArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistService := mocks.NewMockArtistServiceI(ctrl)

	handler := &ArtistHandler{
		ArtistService: mockArtistService,
	}

	tests := []struct {
		name           string
		method         string
		body           any
		role           string
		mockService    func()
		expectedStatus int
	}{
		{
			name:   "admin creates artist successfully",
			method: http.MethodPost,
			body: map[string]string{
				"name": "Test Artist",
				"bio":  "A very talented artist",
			},
			role: "Admin",
			mockService: func() {
				mockArtistService.EXPECT().CreateArtist(gomock.Any(), "Test Artist", "A very talented artist").
					Return(models.Artist{ID: "123", Name: "Test Artist", Bio: "A very talented artist"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "non admin role-failure",
			method: http.MethodPost,
			body: map[string]string{
				"name": "Test Artist",
				"bio":  "A very talented artist",
			},
			role:           "Customer",
			mockService:    func() {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "wrong http method",
			method:         http.MethodGet,
			body:           nil,
			role:           "Admin",
			mockService:    func() {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid request body",
			method:         http.MethodPost,
			body:           "invalid-json",
			role:           "Admin",
			mockService:    func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "no artist name",
			method: http.MethodPost,
			body: map[string]string{
				"name": "",
				"bio":  "A bio",
			},
			role:           "Admin",
			mockService:    func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "service returns error",
			method: http.MethodPost,
			body: map[string]string{
				"name": "Artist",
				"bio":  "Too short",
			},
			role: "Admin",
			mockService: func() {
				mockArtistService.EXPECT().CreateArtist(gomock.Any(), "Artist", "Too short").
					Return(models.Artist{}, errors.New("bio too short"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if s, ok := tt.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/artists", bytes.NewReader(reqBody))
			req = req.WithContext(AddUserToContext(req.Context(), tt.role))

			rr := httptest.NewRecorder()

			tt.mockService()

			handler.CreateArtist(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func AddUserToContext(ctx context.Context, role string) context.Context {
	ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, role)
	return ctx
}
func TestArtistHandler_DeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistService := mocks.NewMockArtistServiceI(ctrl)

	handler := &ArtistHandler{
		ArtistService: mockArtistService,
	}

	tests := []struct {
		name           string
		method         string
		role           string
		artistID       string
		mockService    func()
		expectedStatus int
	}{
		{
			name:     "admin deletes artist successfully",
			method:   http.MethodDelete,
			role:     "Admin",
			artistID: "123",
			mockService: func() {
				mockArtistService.EXPECT().DeleteArtist(gomock.Any(), "123").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "non admin role - forbidden",
			method:         http.MethodDelete,
			role:           "Customer",
			artistID:       "123",
			mockService:    func() {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "wrong http method",
			method:         http.MethodGet,
			role:           "Admin",
			artistID:       "123",
			mockService:    func() {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "missing artist ID",
			method:         http.MethodDelete,
			role:           "Admin",
			artistID:       "",
			mockService:    func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service returns error",
			method:   http.MethodDelete,
			role:     "Admin",
			artistID: "456",
			mockService: func() {
				mockArtistService.EXPECT().DeleteArtist(gomock.Any(), "456").
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()

			req := httptest.NewRequest(tt.method, "/artists/"+tt.artistID, nil)
			req = req.WithContext(AddUserToContext(req.Context(), tt.role))

			req.SetPathValue("artistID", tt.artistID)

			rr := httptest.NewRecorder()
			handler.DeleteArtist(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestArtistHandler_BrowseArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistService := mocks.NewMockArtistServiceI(ctrl)

	handler := &ArtistHandler{
		ArtistService: mockArtistService,
	}

	tests := []struct {
		name           string
		method         string
		query          string
		mockService    func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "browse by name - success",
			method: http.MethodGet,
			query:  "?name=John",
			mockService: func() {
				mockArtistService.EXPECT().GetArtists(gomock.Any(), "John").
					Return([]models.Artist{
						{ID: "1", Name: "John", Bio: "Singer"},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1","name":"John","bio":"Singer"}]`,
		},
		{
			name:   "browse by name - service error",
			method: http.MethodGet,
			query:  "?name=John",
			mockService: func() {
				mockArtistService.EXPECT().GetArtists(gomock.Any(), "John").
					Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to fetch artists\n",
		},
		{
			name:   "browse by artistID - success",
			method: http.MethodGet,
			query:  "?artistID=123",
			mockService: func() {
				mockArtistService.EXPECT().GetArtistByID(gomock.Any(), "123").
					Return(&models.Artist{ID: "123", Name: "Artist One", Bio: "Bio One"}, nil)

			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"123","name":"Artist One","bio":"Bio One"}`,
		},
		{
			name:   "browse by artistID - service error",
			method: http.MethodGet,
			query:  "?artistID=123",
			mockService: func() {
				mockArtistService.EXPECT().GetArtistByID(gomock.Any(), "123").
					Return(nil, errors.New("db error"))

			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to fetch artist\n",
		},
		{
			name:           "wrong method",
			method:         http.MethodPost,
			query:          "",
			mockService:    func() {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
		{
			name:   "no query params",
			method: http.MethodGet,
			query:  "",
			mockService: func() {
				mockArtistService.EXPECT().GetArtistByID(gomock.Any(), "").
					Return(&models.Artist{}, errors.New("not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to fetch artist\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()

			req := httptest.NewRequest(tt.method, "/artists"+tt.query, nil)
			rr := httptest.NewRecorder()

			handler.BrowseArtists(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
