package models

import "time"

type StatusType string

const (
	Active   StatusType = "active"
	Canceled StatusType = "canceled"
)

type Trip struct {
	ID          string     `json:"id"`
	BusID       int        `json:"busID" db:"busID"`
	Origin      string     `json:"origin" db:"origin"`
	Destination string     `json:"destination" db:"destination"`
	TripDate    string     `json:"tripDate" db:"tripDate"`
	Fare        int        `json:"fare" db:"fare"`
	Status      StatusType `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
