package models

type UpdateUserRequest struct {
	IsBlocked *bool   `json:"isBlocked,omitempty"`
	Role      *string `json:"role,omitempty"`
}
