package authorisation

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	userID := "user123"
	email := "user@example.com"
	role := "admin"

	tokenString, err := GenerateJWT(userID, email, role)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Claims.UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Claims.Email = %v, want %v", claims.Email, email)
	}
	if claims.Role != role {
		t.Errorf("Claims.Role = %v, want %v", claims.Role, role)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Errorf("Token expiry time is in the past")
	}
}

func TestValidateJWT(t *testing.T) {
	validToken, err := GenerateJWT("user123", "user@example.com", "admin")
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
		verify  func(claims *Claims) error
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
			verify: func(claims *Claims) error {
				if claims.UserID != "user123" {
					return fmt.Errorf("UserID = %v, want %v", claims.UserID, "user123")
				}
				if claims.Email != "user@example.com" {
					return fmt.Errorf("Email = %v, want %v", claims.Email, "user@example.com")
				}
				if claims.Role != "admin" {
					return fmt.Errorf("Role = %v, want %v", claims.Role, "admin")
				}
				if claims.ExpiresAt.Time.Before(time.Now()) {
					return fmt.Errorf("token expired")
				}
				return nil
			},
		},
		{
			name:    "invalid token string",
			token:   "invalid.token.string",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.verify != nil {
				if err := tt.verify(got); err != nil {
					t.Errorf("ValidateJWT() claims verification failed: %v", err)
				}
			}
		})
	}
}
