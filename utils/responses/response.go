package responses

import (
	"encoding/json"
	"net/http"

	"eventro2/models"
)

func InvalidRequest(w http.ResponseWriter) {
	resp := models.CustomResponse{
		Message:    "Invalid request body",
		StatusCode: http.StatusBadRequest,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func UnauthorisedRequest(w http.ResponseWriter) {
	resp := models.CustomResponse{
		Message:    "Unauthorized",
		StatusCode: http.StatusUnauthorized,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func MethodNotAllowed(w http.ResponseWriter) {
	resp := models.CustomResponse{
		Message:    "Method not allowed",
		StatusCode: http.StatusMethodNotAllowed,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func InternalServerError(w http.ResponseWriter, message string) {
	resp := models.CustomResponse{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func CustomError(w http.ResponseWriter, statusCode int, Message string) {
	resp := models.CustomResponse{
		Message:    Message,
		StatusCode: statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func Forbidden(w http.ResponseWriter) {
	resp := models.CustomResponse{
		Message:    "Forbidden",
		StatusCode: http.StatusForbidden,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
