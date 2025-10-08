package handlers

import (
	"encoding/json"
	"eventro2/middleware"
	"eventro2/models"
	"eventro2/services/userservice"
	"eventro2/utils/responses"
	"net/http"
	"strings"
)

type UserHandler struct {
	UserService userservice.UserServiceI
}

func (h *UserHandler) BrowseUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userRole, err := middleware.GetUserRole(r.Context())
	userRole = strings.ToLower(userRole)
	if err != nil || userRole != "admin" {
		http.Error(w, "Forbidden: admin only", http.StatusForbidden)
		return
	}

	query := r.URL.Query()
	blockedParam := query.Get("blocked")
	userID := query.Get("userId")

	var blocked *bool
	if blockedParam != "" {
		val := blockedParam == "true"
		blocked = &val
	}

	if userID == "" {
		users, err := h.UserService.BrowseUsers(r.Context(), userID, blocked)
		if err != nil {
			http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	} else {
		user, err := h.UserService.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure user is admin
	userRole, err := middleware.GetUserRole(r.Context())
	userRole = strings.ToLower(userRole)
	if err != nil || userRole != "admin" {
		http.Error(w, "Forbidden: admin only", http.StatusForbidden)
		return
	}

	// Extract userId from URL
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update user
	updatedUser, err := h.UserService.UpdateUser(r.Context(), userID, req)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) GetUserByMailID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	mailID := r.PathValue("mailID")
	if mailID == "" {
		responses.InvalidRequest(w)
		return
	}

	user, err := h.UserService.GetUserByMailID(r.Context(), mailID)
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.PathValue("userID")
	if userID == "" {
		responses.InvalidRequest(w)
		return
	}

	user, err := h.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}
