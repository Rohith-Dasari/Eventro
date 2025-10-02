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

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		role           string
		body           any
		mockSetup      func(m *mocks.MockUserServiceI)
		expectedStatus int
	}{
		{
			name:   "invalid method",
			method: http.MethodGet,
			userID: "123",
			role:   "Admin",
			body:   models.UpdateUserRequest{},
			mockSetup: func(m *mocks.MockUserServiceI) {
				// no calls expected
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "missing userID",
			method: http.MethodPatch,
			userID: "",
			role:   "Admin",
			body:   models.UpdateUserRequest{},
			mockSetup: func(m *mocks.MockUserServiceI) {
				// no calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid JSON body",
			method: http.MethodPatch,
			userID: "123",
			role:   "Admin",
			body:   "not-json",
			mockSetup: func(m *mocks.MockUserServiceI) {
				// no calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "forbidden - non-admin user",
			method: http.MethodPatch,
			userID: "123",
			role:   "Customer",
			body:   models.UpdateUserRequest{Role: stringPtr("Admin")},
			mockSetup: func(m *mocks.MockUserServiceI) {
				// no calls expected
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "service error",
			method: http.MethodPatch,
			userID: "123",
			role:   "Admin",
			body:   models.UpdateUserRequest{IsBlocked: boolPtr(true)},
			mockSetup: func(m *mocks.MockUserServiceI) {
				m.EXPECT().
					UpdateUser(gomock.Any(), "123", gomock.Any()).
					Return(models.User{}, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "success - update user",
			method: http.MethodPatch,
			userID: "123",
			role:   "Admin",
			body:   models.UpdateUserRequest{Role: stringPtr("Customer"), IsBlocked: boolPtr(false)},
			mockSetup: func(m *mocks.MockUserServiceI) {
				m.EXPECT().
					UpdateUser(gomock.Any(), "123", gomock.Any()).
					Return(models.User{UserID: "123", Role: "Customer", IsBlocked: false}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocks.NewMockUserServiceI(ctrl)
			tt.mockSetup(mockUserService)

			h := &UserHandler{UserService: mockUserService}

			var bodyBytes []byte
			var err error
			if str, ok := tt.body.(string); ok {
				// invalid JSON case
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/users/"+tt.userID, bytes.NewReader(bodyBytes))
			if tt.userID != "" {
				req.SetPathValue("userID", tt.userID)
			}

			ctx := context.WithValue(req.Context(), middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdateUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }

func TestUserHandler_BrowseUsers(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		role           string
		query          string
		mockSetup      func(m *mocks.MockUserServiceI)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "invalid method",
			method: http.MethodPost,
			role:   "Admin",
			query:  "",
			mockSetup: func(m *mocks.MockUserServiceI) {
				// no calls expected
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "forbidden - not admin",
			method: http.MethodGet,
			role:   "Customer",
			query:  "",
			mockSetup: func(m *mocks.MockUserServiceI) {
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "success - browse all users",
			method: http.MethodGet,
			role:   "Admin",
			query:  "",
			mockSetup: func(m *mocks.MockUserServiceI) {
				m.EXPECT().
					BrowseUsers(gomock.Any(), "", (*bool)(nil)).
					Return([]models.User{
						{UserID: "u1", Role: "Host"},
						{UserID: "u2", Role: "Customer"},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []models.User{
				{UserID: "u1", Role: "Host"},
				{UserID: "u2", Role: "Customer"},
			},
		},
		{
			name:   "success - browse by userId",
			method: http.MethodGet,
			role:   "Admin",
			query:  "?userId=u1",
			mockSetup: func(m *mocks.MockUserServiceI) {
				m.EXPECT().
					GetUserByID(gomock.Any(), "u1").
					Return(&models.User{UserID: "u1", Role: "Host"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   &models.User{UserID: "u1", Role: "Host"},
		},
		{
			name:   "service error - browse users",
			method: http.MethodGet,
			role:   "Admin",
			query:  "",
			mockSetup: func(m *mocks.MockUserServiceI) {
				m.EXPECT().
					BrowseUsers(gomock.Any(), "", (*bool)(nil)).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "service error - get user by ID",
			method: http.MethodGet,
			role:   "Admin",
			query:  "?userId=u1",
			mockSetup: func(m *mocks.MockUserServiceI) {
				m.EXPECT().
					GetUserByID(gomock.Any(), "u1").
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockUserServiceI(ctrl)
			tt.mockSetup(mockSvc)

			handler := &UserHandler{UserService: mockSvc}

			req := httptest.NewRequest(tt.method, "/users"+tt.query, nil)
			ctx := context.WithValue(req.Context(), middleware.ContextUserRoleKey, tt.role)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.BrowseUsers(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != nil && rr.Code == http.StatusOK {
				if users, ok := tt.expectedBody.([]models.User); ok {
					var got []models.User
					if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
						t.Fatalf("failed to decode body: %v", err)
					}
					assert.Equal(t, users, got)
				} else if user, ok := tt.expectedBody.(*models.User); ok {
					var got models.User
					if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
						t.Fatalf("failed to decode body: %v", err)
					}
					assert.Equal(t, *user, got)
				}
			}
		})
	}
}
