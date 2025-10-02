package handlers

import (
	"eventro2/mocks"
	"eventro2/models"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := &AuthHandler{AuthService: mockAuthService}

	tests := []struct {
		name           string
		setupMock      func()
		requestBody    string
		expectedStatus int
		expectToken    bool
		method         string
	}{
		{
			name: "success",
			setupMock: func() {
				mockAuthService.EXPECT().
					ValidateLogin(gomock.Any(), "test@example.com", "password123").
					Return(models.User{
						UserID:   "123",
						Username: "testuser",
						Email:    "test@example.com",
						Role:     "user",
					}, nil)
			},
			requestBody:    `{"email":"test@example.com","password":"password123"}`,
			expectedStatus: http.StatusOK,
			expectToken:    true,
			method:         http.MethodPost,
		},
		{
			name:           "invalid json",
			setupMock:      func() {},       // no calls expected
			requestBody:    `{"email":123}`, // bad type
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
			method:         http.MethodPost,
		},
		{
			name: "invalid credentials",
			setupMock: func() {
				mockAuthService.EXPECT().
					ValidateLogin(gomock.Any(), "wrong@example.com", "badpass").
					Return(models.User{}, fmt.Errorf("invalid credentials"))
			},
			requestBody:    `{"email":"wrong@example.com","password":"badpass"}`,
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
			method:         http.MethodPost,
		},
		{
			name:           "method not allowed",
			setupMock:      func() {},
			requestBody:    "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectToken:    false,
			method:         http.MethodGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)

			if tt.expectToken {
				if !strings.Contains(string(body), "token") {
					t.Errorf("expected token in response, got %s", body)
				}
			}
		})
	}
}

func TestAuthHandler_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := &AuthHandler{AuthService: mockAuthService}

	tests := []struct {
		name           string
		setupMock      func()
		requestBody    string
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockAuthService.EXPECT().
					Signup(gomock.Any(), "testuser", "test@example.com", "1234567890", "password123").
					Return(models.User{
						UserID:   "123",
						Username: "testuser",
						Email:    "test@example.com",
						Role:     "user",
					}, nil)
			},
			requestBody:    `{"username":"testuser","email":"test@example.com","phone_number":"1234567890","password":"password123"}`,
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name:           "invalid json",
			setupMock:      func() {},
			requestBody:    `{"username":123}`,
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name: "signup error",
			setupMock: func() {
				mockAuthService.EXPECT().
					Signup(gomock.Any(), "baduser", "bad@example.com", "000", "nopass").
					Return(models.User{}, fmt.Errorf("signup failed"))
			},
			requestBody:    `{"username":"baduser","email":"bad@example.com","phone_number":"000","password":"nopass"}`,
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			handler.Signup(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)

			if tt.expectToken {
				if !strings.Contains(string(body), "token") {
					t.Errorf("expected token in response, got %s", body)
				}
			}
		})
	}
}
