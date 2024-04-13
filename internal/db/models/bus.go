package models

import "time"

type BusType string

const (
	Minibus  BusType = "minibus"
	Midibus  BusType = "midibus"
	StaffBus BusType = "staffbus"
)

type Bus struct {
	PlateNumber  string    `json:"plate_number" db:"pate_number"`
	Manufacturer string    `json:"manufacturer" db:"manufacturer"`
	Model        string    `json:"model" db:"model"`
	Type         BusType   `json:"type" db:"type"`
	Year         string    `json:"year" db:"year"`
	Capacity     int       `json:"capacity" db:"capacity"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type BusReqBody struct {
	PlateNumber  string  `json:"plate_number" db:"plate_number"`
	Manufacturer string  `json:"manufacturer" db:"manufacturer"`
	Model        string  `json:"model" db:"model"`
	Type         BusType `json:"type" db:"type"`
	Year         string  `json:"year" db:"year"`
	Capacity     int     `json:"capacity" db:"capacity"`
}

type CreatedBus struct {
	ID           int       `json:"id" db:"id"`
	PlateNumber  string    `json:"plate_number" db:"pate_number"`
	Manufacturer string    `json:"manufacturer" db:"manufacturer"`
	Model        string    `json:"model" db:"model"`
	Type         BusType   `json:"type" db:"type"`
	Year         string    `json:"year" db:"year"`
	Capacity     int       `json:"capacity" db:"capacity"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
