package handlers

import (
	"encoding/json"
	"eventro2/middleware"
	"eventro2/models"
	"eventro2/services/eventservice"
	"eventro2/utils/responses"
	"fmt"
	"net/http"
)

type EventHandler struct {
	EventService eventservice.EventServiceI
}

type CreateEventRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	Category    string   `json:"category"`
	Artists     []string `json:"artists"`
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		responses.MethodNotAllowed(w)
		return
	}

	// Admin authorization check
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "Admin" {
		responses.Forbidden(w)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	event, err := h.EventService.CreateNewEvent(
		r.Context(),
		req.Name,
		req.Description,
		req.Duration,
		models.EventCategory(req.Category),
		req.Artists,
	)
	if err != nil {
		responses.InternalServerError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) BrowseEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit")
	if r.Method != http.MethodGet {
		responses.MethodNotAllowed(w)
		return
	}

	query := r.URL.Query()
	isBlockedParam := query.Get("isBlocked")

	var isBlocked *bool
	if isBlockedParam != "" {
		role, err := middleware.GetUserRole(r.Context())
		if err != nil || role != "Admin" {
			responses.Forbidden(w)
			return
		}
		//change
		val := isBlockedParam == "true"
		isBlocked = &val
	}

	filter := models.EventFilter{
		EventID:    query.Get("eventID"),
		Name:       query.Get("eventname"),
		Category:   query.Get("category"),
		Location:   query.Get("location"),
		IsBlocked:  isBlocked,
		ArtistName: query.Get("artistName"),
	}

	events, err := h.EventService.BrowseEvents(r.Context(), filter)
	if err != nil {
		responses.InternalServerError(w, "Failed to fetch events")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responses.MethodNotAllowed(w)
		return
	}

	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "Admin" {
		responses.Forbidden(w)
		return
	}

	eventID := r.PathValue("eventID")
	if eventID == "" {
		responses.InvalidRequest(w)
		return
	}

	err = h.EventService.DeleteEvent(r.Context(), eventID)
	if err != nil {
		responses.InternalServerError(w, "Failed to delete event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		responses.MethodNotAllowed(w)
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "Admin" {
		responses.Forbidden(w)
		return
	}
	//change

	eventID := r.PathValue("eventID")
	if eventID == "" {
		responses.InvalidRequest(w)
		return
	}

	var updateData models.EventUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		responses.InvalidRequest(w)
		return
	}

	updatedEvent, err := h.EventService.UpdateEvent(r.Context(), eventID, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEvent)
}
