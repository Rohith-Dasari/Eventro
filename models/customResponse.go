package models

type CustomResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
