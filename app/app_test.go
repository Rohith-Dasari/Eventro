package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"
)

func TestSetupServer(t *testing.T) {

	mockDB := &gorm.DB{}

	mux := SetupServer(mockDB)
	if mux == nil {
		t.Fatal("SetupServer returned nil mux")
	}

	testEndpoints := []struct {
		method string
		path   string
		auth   bool
	}{
		{"POST", "/api/v1/login", false},
		{"POST", "/api/v1/signup", false},
		{"GET", "/api/v1/artists", true},
		{"POST", "/api/v1/events", true},
		{"GET", "/api/v1/shows", true},
		{"POST", "/api/v1/venues", true},
		{"GET", "/api/v1/bookings", true},
	}

	for _, endpoint := range testEndpoints {
		req := httptest.NewRequest(endpoint.method, endpoint.path, nil)

		if endpoint.auth {
			req.Header.Set("Authorization", "Bearer test-token")
			continue
		}

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code == http.StatusNotFound {
			t.Errorf("Route %s %s not found", endpoint.method, endpoint.path)
		}
	}
}
