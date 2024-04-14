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
	TripDate    string     `json:"trip_date" db:"trip_date"`
	Fare        int        `json:"fare" db:"fare"`
	Status      StatusType `json:"status" db:"status"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type TripReqBody struct {
	BusID       int        `json:"busID" db:"busID"`
	Origin      string     `json:"origin" db:"origin"`
	Destination string     `json:"destination" db:"destination"`
	TripDate    string     `json:"trip_date" db:"trip_date"`
	Fare        int        `json:"fare" db:"fare"`
	Status      StatusType `json:"status" db:"status"`
}

type CancelTrip struct {
	Status    StatusType `json:"status" db:"status"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type CancelTripReqBody struct {
	Status StatusType `json:"status" db:"status"`
}
