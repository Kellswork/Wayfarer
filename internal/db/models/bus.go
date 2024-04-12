package models

import "time"

type Bus struct {
	PlateNumber  string    `json:"plate_number" db:"pate_number"`
	Manufacturer string    `json:"manufacturer" db:"manufacturer"`
	Model        string    `json:"model" db:"model"`
	Type         BusType   `json:"type" db:"type"`
	Year         string    `json:"year" db:"year"`
	Capacity     string    `json:"capacity" db:"capacity"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type BusType struct {
	Minibus  string `json:"minibus"`
	Midibus  string `json:"midibus"`
	StaffBus string `json:"staffbus"`
}

type BusReqBody struct {
	PlateNumber  string  `json:"plate_number" db:"pate_number"`
	Manufacturer string  `json:"manufacturer" db:"manufacturer"`
	Model        string  `json:"model" db:"model"`
	Type         BusType `json:"type" db:"type"`
	Year         string  `json:"year" db:"year"`
}
