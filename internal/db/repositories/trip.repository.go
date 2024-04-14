package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/kellswork/wayfarer/internal/db/models"
)

//go:generate /Users/kells/go/bin/mockgen -source trip_repository.go -destination ./mocks/trip_repository.go -package mocks repositories TripRepository

type TripRepository interface {
	Create(ctx context.Context, bus *models.Trip) error
	Cancel(ctx context.Context, status string, updatedAt time.Time, tripID string) error
	FindAll(ctx context.Context) (*[]models.Trip, error)
}

type tripRepository struct {
	db *sql.DB
}

func newTripRepository(db *sql.DB) *tripRepository {
	return &tripRepository{
		db: db,
	}
}

func (tr *tripRepository) Create(ctx context.Context, trip *models.Trip) error {

	query := `INSERT INTO trips ("id", "busID", "origin", "destination", "trip_date", "fare", "status", "created_at") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tr.db.ExecContext(ctx, query, trip.ID, trip.BusID, trip.Origin, trip.Destination, trip.TripDate, trip.Fare, trip.Status, trip.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (tr *tripRepository) Cancel(ctx context.Context, status string, updatedAt time.Time, tripID string) error {

	query := `UPDATE trips SET status=$1, updated_at=$2 WHERE id =$3`

	_, err := tr.db.ExecContext(ctx, query, status, updatedAt, tripID)

	if err != nil {
		return err
	}

	return nil

	// .Scan(&trip.ID, &trip.BusID, &trip.Origin, &trip.Destination, &trip.TripDate, &trip.Fare, &trip.Status, &trip.CreatedAt, &trip.UpdatedAt)
}

func (tr *tripRepository) FindAll(ctx context.Context) (*[]models.Trip, error) {

	query := "SELECT * FROM trips"

	rows, err := tr.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := make([]models.Trip, 0)
	for rows.Next() {
		var trip models.Trip
		if err := rows.Scan(&trip.ID, &trip.BusID, &trip.Origin, &trip.Destination, &trip.TripDate, &trip.Fare, &trip.Status, &trip.CreatedAt, &trip.UpdatedAt); err != nil {
			return nil, err
		}
		trips = append(trips, trip)

		// check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return &trips, nil
}
