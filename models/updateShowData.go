package models

import "time"

type UpdateShowData struct {
	Price     *float64
	ShowDate  *time.Time
	ShowTime  *string
	IsBlocked *bool
}
