package handlers

import (
	"encoding/json"
	"eventro2/middleware"
	"eventro2/models"
	"eventro2/services/venueservice"
	"eventro2/utils/responses"
	"fmt"
	"net/http"
	"strings"
)

type VenueHandler struct {
	VenueService venueservice.VenueServiceInterface
}

type CreateVenueRequest struct {
	Name                 string `json:"name"`
	City                 string `json:"city"`
	State                string `json:"state"`
	IsSeatLayoutRequired bool   `json:"is_seat_layout_required"`
}

type UpdateVenueRequest struct {
	Name                 *string `json:"name,omitempty"`
	City                 *string `json:"city,omitempty"`
	State                *string `json:"state,omitempty"`
	IsSeatLayoutRequired *bool   `json:"is_seat_layout_required,omitempty"`
	IsBlocked            *bool   `json:"is_blocked,omitempty"`
}

func (h *VenueHandler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.MethodNotAllowed(w)
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "Host" {
		responses.Forbidden(w)
		//change
		return
	}
	hostID, err := middleware.GetUserID(r.Context())
	if err != nil || hostID == "" {
		responses.UnauthorisedRequest(w)
		return
	}

	var req CreateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	venue, err := h.VenueService.CreateVenue(
		r.Context(),
		hostID,
		req.Name,
		req.City,
		req.State,
		req.IsSeatLayoutRequired,
	)
	if err != nil {
		responses.InternalServerError(w, "Failed to create venue: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(venue)
}
func (h *VenueHandler) UpdateVenue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		responses.MethodNotAllowed(w)
		return
	}

	venueID := r.PathValue("venueID")
	if venueID == "" {
		responses.InvalidRequest(w)
		return
	}
	fmt.Println("reaching till hereeeee")

	userID, err := middleware.GetUserID(r.Context())
	if err != nil || userID == "" {
		responses.UnauthorisedRequest(w)
		return
	}
	fmt.Println("reaching till here")
	userRole, err := middleware.GetUserRole(r.Context())
	if err != nil || userRole != "Host" {
		responses.UnauthorisedRequest(w)
		return
	}
	fmt.Println("reaching till here with role")

	var req UpdateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}
	fmt.Println("decoded")

	update := models.UpdateVenueData{
		Name:                 req.Name,
		City:                 req.City,
		State:                req.State,
		IsSeatLayoutRequired: req.IsSeatLayoutRequired,
		IsBlocked:            req.IsBlocked,
	}

	updatedVenue, err := h.VenueService.UpdateVenue(r.Context(), venueID, userID, userRole, update)
	if err != nil {
		fmt.Println(err)
		responses.Forbidden(w)
		return
	}
	fmt.Println("done")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedVenue)
}

func (h *VenueHandler) DeleteVenue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responses.MethodNotAllowed(w)
		return
	}

	venueID := r.PathValue("venueID")
	if venueID == "" {
		responses.InvalidRequest(w)
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil || userID == "" {
		responses.UnauthorisedRequest(w)
		return
	}

	userRole, err := middleware.GetUserRole(r.Context())
	userRole = strings.ToLower(userRole)
	if err != nil || (userRole != "host" && userRole != "admin") {
		responses.Forbidden(w)
		return
	}

	if err := h.VenueService.DeleteVenue(r.Context(), venueID, userID, userRole); err != nil {
		responses.Forbidden(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *VenueHandler) BrowseVenues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.MethodNotAllowed(w)
		return
	}
	if r.URL.Query().Get("hostId") != "" || r.URL.Query().Get("isBlocked") != "" {
		userRole, err := middleware.GetUserRole(r.Context())
		userRole = strings.ToLower(userRole)
		if err != nil || (userRole != "host" && userRole != "admin") {
			responses.UnauthorisedRequest(w)
			return
		}
	}
	//change

	query := r.URL.Query()
	filter := models.VenueFilter{
		City:      query.Get("city"),
		HostID:    query.Get("hostId"),
		VenueID:   query.Get("venueId"),
		IsBlocked: query.Get("isBlocked") == "true",
	}

	venues, err := h.VenueService.BrowseVenues(r.Context(), filter)
	if err != nil {
		responses.InternalServerError(w, "Failed to fetch venues: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(venues)
}
