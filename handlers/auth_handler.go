package handlers

import (
	"encoding/json"
	"eventro2/services/authorisation"
	"eventro2/utils/responses"
	"net/http"
)

type AuthHandler struct {
	AuthService authorisation.AuthServiceInterface
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type SignupRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.MethodNotAllowed(w)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}

	user, err := h.AuthService.ValidateLogin(r.Context(), req.Email, req.Password)
	if err != nil {
		responses.UnauthorisedRequest(w)
		return
	}

	token, err := authorisation.GenerateJWT(user.UserID, user.Email, string(user.Role))
	if err != nil {
		responses.InternalServerError(w, "Failed to generate token")
		return
	}

	res := LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		responses.InternalServerError(w, "Failed to encode response")
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.MethodNotAllowed(w)
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.InvalidRequest(w)
		return
	}
	if req.Username == "" || req.Email == "" || req.PhoneNumber == "" || req.Password == "" {
		responses.InvalidRequest(w)
		return
	}
	if len(req.Password) < 12 {
		responses.CustomError(w, 400, "Password should be of atleast 12 alphanumeric characters and a symbol")
		return
	}

	user, err := h.AuthService.Signup(r.Context(), req.Username, req.Email, req.PhoneNumber, req.Password)
	if err != nil {
		responses.CustomError(w, 409, "user already exists")
		return
		///change
	}

	token, err := authorisation.GenerateJWT(user.UserID, user.Email, string(user.Role))
	if err != nil {
		responses.InternalServerError(w, "Failed to generate token")
		return
	}

	res := SignupResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		responses.InternalServerError(w, "Failed to encode response")
	}
}
