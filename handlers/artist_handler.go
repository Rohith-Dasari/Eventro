package handlers

import (
	"encoding/json"
	"eventro2/middleware"
	"eventro2/services/artistservice"
	"eventro2/utils/responses"
	"log"
	"net/http"
)

type CreateArtistRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type CreateArtistResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type ArtistHandler struct {
	ArtistService artistservice.ArtistServiceI
}

func (h *ArtistHandler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit CreateArtist endpoint")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Admin authorization check
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "Admin" {
		http.Error(w, "Forbidden: admin only", http.StatusForbidden)
		return
	}

	var req CreateArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	if req.Name == "" {
		http.Error(w, "Artist name is required", http.StatusBadRequest)
		return
	}

	artist, err := h.ArtistService.CreateArtist(r.Context(), req.Name, req.Bio)
	if err != nil {
		http.Error(w, "Failed to create artist", http.StatusInternalServerError)
		return
	}

	res := CreateArtistResponse{
		ID:   artist.ID,
		Name: artist.Name,
		Bio:  artist.Bio,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *ArtistHandler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "Admin" {
		http.Error(w, "Forbidden: admin only", http.StatusForbidden)
		return
	}
	id := r.PathValue("artistID")
	if id == "" {
		http.Error(w, "Artist ID required", http.StatusBadRequest)
		return
	}

	err = h.ArtistService.DeleteArtist(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete artist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ArtistHandler) BrowseArtists(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit findArtist endpoint")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	artistID := r.URL.Query().Get("artistID")
	if artistID != "" {
		artist, err := h.ArtistService.GetArtistByID(r.Context(), artistID)
		if err != nil {
			http.Error(w, "Failed to fetch artist", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(artist)
	} else {
		artists, err := h.ArtistService.GetArtists(r.Context(), name)
		if err != nil {
			http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(artists)
	}
}

//change
