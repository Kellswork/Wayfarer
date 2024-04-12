package repositories

import (
	"context"
	"database/sql"

	"github.com/kellswork/wayfarer/internal/db/models"
)

//go:generate /Users/kells/go/bin/mockgen -source user_repository.go -destination ./mocks/user_repository.go -package mocks repositories UserRepository

type BusRepository interface {
	Create(ctx context.Context, user *models.Bus) error
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
