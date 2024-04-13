package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kellswork/wayfarer/internal/db/models"
)

//go:generate /Users/kells/go/bin/mockgen -source user_repository.go -destination ./mocks/user_repository.go -package mocks repositories UserRepository

type BusRepository interface {
	Create(ctx context.Context, bus *models.Bus) error
	DoesPlateExists(ctx context.Context, plateNumber string) (bool, error)
}

type busRespository struct {
	db *sql.DB
}

func newBusRepository(db *sql.DB) *busRespository {
	return &busRespository{
		db: db,
	}
}

func (ur *busRespository) Create(ctx context.Context, bus *models.Bus) error {
	query := "INSERT INTO bus ( plate_number, manufacturer, model, type, year, capacity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *"

	_, err := ur.db.Exec(query, bus.PlateNumber, bus.Manufacturer, bus.Model, bus.Type, bus.Year, bus.Capacity, bus.CreatedAt, bus.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (ur *userRespository) DoesPlateExists(ctx context.Context, plateNumber string) (bool, error) {
	var plateCount int

	query := "SELECT COUNT(*) FROM bus WHERE plate_number = $1"
	err := ur.db.QueryRowContext(ctx, query, plateNumber).Scan(&plateCount)
	if err != nil {
		fmt.Printf("checking if email exists: %v\n", err)
		return false, err
	}
	return plateCount > 0, nil
}
