// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package postgresdriver

import (
	"time"
)

type RelayCount struct {
	AppPublicKey string    `json:"appPublicKey"`
	Day          time.Time `json:"day"`
	Success      int32     `json:"success"`
	Error        int32     `json:"error"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}