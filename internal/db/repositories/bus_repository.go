package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kellswork/wayfarer/internal/db/models"
)

//go:generate /Users/kells/go/bin/mockgen -source user_repository.go -destination ./mocks/user_repository.go -package mocks repositories UserRepository

type BusRepository interface {
	Create(ctx context.Context, bus *models.Bus) (*models.CreatedBus, error)
	DoesPlateExists(ctx context.Context, plateNumber string) (bool, error)
	FindAllBuses(ctx context.Context) (*[]models.CreatedBus, error)
}

type busRespository struct {
	db *sql.DB
}

func newBusRepository(db *sql.DB) *busRespository {
	return &busRespository{
		db: db,
	}
}

func (br *busRespository) Create(ctx context.Context, bus *models.Bus) (*models.CreatedBus, error) {
	var r models.CreatedBus

	query := "INSERT INTO bus ( plate_number, manufacturer, model, type, year, capacity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *"

	err := br.db.QueryRowContext(ctx, query, bus.PlateNumber, bus.Manufacturer, bus.Model, bus.Type, bus.Year, bus.Capacity, bus.CreatedAt, bus.UpdatedAt).Scan(&r.ID, &r.PlateNumber, &r.Manufacturer, &r.Model, &r.Type, &r.Year, &r.Capacity, &r.CreatedAt, &r.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (br *busRespository) DoesPlateExists(ctx context.Context, plateNumber string) (bool, error) {
	var plateCount int

	query := "SELECT COUNT(*) FROM bus WHERE plate_number = $1"
	err := br.db.QueryRowContext(ctx, query, plateNumber).Scan(&plateCount)
	if err != nil {
		fmt.Printf("checking if email exists: %v\n", err)
		return false, err
	}
	return plateCount > 0, nil
}

func (br *busRespository) FindAllBuses(ctx context.Context) (*[]models.CreatedBus, error) {

	query := "SELECT * FROM bus"

	rows, err := br.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buses := make([]models.CreatedBus, 0)
	for rows.Next() {
		var bus models.CreatedBus
		if err := rows.Scan(&bus.ID, &bus.PlateNumber, &bus.Manufacturer, &bus.Model, &bus.Type, &bus.Year, &bus.Capacity, &bus.CreatedAt, &bus.UpdatedAt); err != nil {
			return nil, err
		}
		buses = append(buses, bus)

		// check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return &buses, nil

}
