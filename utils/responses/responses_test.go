package responses_test

import (
	"encoding/json"
	"errors"
	"eventro2/models"
	"eventro2/utils/responses"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponses(t *testing.T) {
	tests := []struct {
		name           string
		fn             func(http.ResponseWriter)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "InvalidRequest",
			fn:             responses.InvalidRequest,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Invalid request body",
		},
		{
			name:           "UnauthorisedRequest",
			fn:             responses.UnauthorisedRequest,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "MethodNotAllowed",
			fn:             responses.MethodNotAllowed,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedMsg:    "Method not allowed",
		},
		{
			name: "InternalServerError",
			fn: func(w http.ResponseWriter) {
				responses.InternalServerError(w, "internal error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "internal error",
		},
		{
			name: "CustomError",
			fn: func(w http.ResponseWriter) {
				responses.CustomError(w, http.StatusTeapot, "custom message")
			},
			expectedStatus: http.StatusTeapot,
			expectedMsg:    "custom message",
		},
		{
			name:           "Forbidden",
			fn:             responses.Forbidden,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			tt.fn(w)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			var body models.CustomResponse
			err := json.NewDecoder(res.Body).Decode(&body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, body.Message)
			assert.Equal(t, tt.expectedStatus, body.StatusCode)
		})
	}
}

type failingWriter struct{}

func (f *failingWriter) Header() http.Header        { return http.Header{} }
func (f *failingWriter) Write([]byte) (int, error)  { return 0, errors.New("write failed") }
func (f *failingWriter) WriteHeader(statusCode int) {}

func TestResponses_EncodingError(t *testing.T) {
	funcs := []func(http.ResponseWriter){
		responses.InvalidRequest,
		responses.UnauthorisedRequest,
		responses.MethodNotAllowed,
		func(w http.ResponseWriter) { responses.InternalServerError(w, "error") },
		func(w http.ResponseWriter) { responses.CustomError(w, 418, "error") },
		responses.Forbidden,
	}

	for _, fn := range funcs {
		fn(&failingWriter{}) // force encoding error
	}
}
