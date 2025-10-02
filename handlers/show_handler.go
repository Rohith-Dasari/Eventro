package handlers

import (
	"encoding/json"
	"eventro2/middleware"
	"eventro2/models"
	"eventro2/services/showservice"
	"eventro2/utils/responses"
	"net/http"
	"time"
)

type ShowHandler struct {
	ShowService showservice.ShowServiceInterface
}

type CreateShowRequest struct {
	EventID  string  `json:"event_id"`
	VenueID  string  `json:"venue_id"`
	Price    float64 `json:"price"`
	ShowDate string  `json:"show_date"`
	ShowTime string  `json:"show_time"`
}

type UpdateShowRequest struct {
	Price     *float64 `json:"price,omitempty"`
	ShowDate  *string  `json:"show_date,omitempty"`
	ShowTime  *string  `json:"show_time,omitempty"`
	IsBlocked *bool    `json:"is_blocked,omitempty"`
}

func (h *ShowHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.MethodNotAllowed(w)
		return
	}

	// get host id
	userID, err := middleware.GetUserID(r.Context())
	if err != nil || userID == "" {
		responses.UnauthorisedRequest(w)
		return
	}
	//get role and if not host say unauthorised
	userRole, err := middleware.GetUserRole(r.Context())
	if err != nil || userRole != "Host" {
		responses.UnauthorisedRequest(w)
		return
	}

	var req CreateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	// Parse ShowDate
	parsedDate, err := time.Parse("2006-01-02", req.ShowDate)
	if err != nil {
		responses.CustomError(w, http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD")
		return
	}

	show, err := h.ShowService.CreateShow(
		r.Context(),
		req.EventID,
		req.VenueID,
		userID,
		req.Price,
		parsedDate,
		req.ShowTime,
	)
	if err != nil {
		responses.CustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(show)
}

func (h *ShowHandler) UpdateShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		responses.MethodNotAllowed(w)
		return
	}

	showID := r.PathValue("showID")

	userID, err := middleware.GetUserID(r.Context())
	if err != nil || userID == "" {
		responses.UnauthorisedRequest(w)
		return
	}

	var req UpdateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	var parsedDate *time.Time
	if req.ShowDate != nil {
		t, err := time.Parse("2006-01-02", *req.ShowDate)
		if err != nil {
			responses.CustomError(w, http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD")
			return
		}
		parsedDate = &t
	}

	update := models.UpdateShowData{
		Price:     req.Price,
		ShowDate:  parsedDate,
		ShowTime:  req.ShowTime,
		IsBlocked: req.IsBlocked,
	}

	updatedShow, err := h.ShowService.UpdateShow(r.Context(), showID, userID, update)
	if err != nil {
		responses.CustomError(w, http.StatusInternalServerError, "Failed to update show: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedShow)
}

func (h *ShowHandler) BrowseShows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.MethodNotAllowed(w)
		return
	}

	query := r.URL.Query()

	filter := models.ShowFilter{
		ShowID:  query.Get("showId"),
		EventID: query.Get("eventId"),
		HostID:  query.Get("hostId"),
		VenueID: query.Get("venueId"),
	}

	shows, err := h.ShowService.BrowseShows(r.Context(), filter)
	if err != nil {
		responses.CustomError(w, http.StatusInternalServerError, "Failed to fetch shows: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shows)
}

func (h *ShowHandler) DeleteShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responses.MethodNotAllowed(w)
		return
	}
	showID := r.PathValue("showID")
	if showID == "" {
		responses.CustomError(w, http.StatusBadRequest, "Show ID required")
		return
	}

	err := h.ShowService.DeleteShow(r.Context(), showID)
	if err != nil {
		responses.CustomError(w, http.StatusInternalServerError, "Failed to delete show: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
