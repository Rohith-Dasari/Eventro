package handlers

import (
	"encoding/json"
	"eventro2/middleware"
	"eventro2/services/bookingservice"
	"eventro2/utils/responses"
	"net/http"
)

// browse bookign, create booking,
type BookingHandler struct {
	BookingService bookingservice.BookingServiceInterface
}

type CreateBookingRequest struct {
	UserID string   `json:"user_id,omitempty"`
	ShowID string   `json:"show_id"`
	Seats  []string `json:"seats"`
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.MethodNotAllowed(w)
		return
	}
	// Only Customer and Admin can create bookings
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || (role != "Customer" && role != "Admin") {
		responses.Forbidden(w)
		return
	}

	authUserID, err := middleware.GetUserID(r.Context())
	if err != nil || authUserID == "" {
		responses.UnauthorisedRequest(w)
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	if req.ShowID == "" || len(req.Seats) == 0 {
		responses.InvalidRequest(w)
		return
	}
	userID := authUserID
	if role == "admin" && req.UserID != "" {
		userID = req.UserID
	}

	booking, err := h.BookingService.AddBooking(r.Context(), userID, req.ShowID, req.Seats)
	if err != nil {
		responses.CustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) BrowseBookings(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodGet {
		responses.MethodNotAllowed(w)
		return
	}

	authUserID, _ := middleware.GetUserID(r.Context())
	role, _ := middleware.GetUserRole(r.Context())

	query := r.URL.Query()
	bookingID := query.Get("bookingId")
	userID := query.Get("userId")
	showID := query.Get("showId")

	if role != "admin" {
		userID = authUserID
	}

	bookings, err := h.BookingService.BrowseBookings(r.Context(), bookingID, userID, showID)
	if err != nil {
		responses.InternalServerError(w, "Failed to fetch bookings")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
