package types

import "time"

type RelayCount struct {
	AppPublicKey string    `json:"appPublicKey"`
	Day          time.Time `json:"day"`
	Success      int       `json:"success"`
	Error        int       `json:"error"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
