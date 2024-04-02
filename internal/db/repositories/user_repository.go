package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kellswork/wayfarer/internal/db/models"
)

//go:generate /Users/kells/go/bin/mockgen -source user_repository.go -destination ./mocks/user_repository.go -package mocks repositories UserRepository

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRespository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) *userRespository {
	return &userRespository{
		db: db,
	}
}

func (ur *userRespository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, first_name, last_name, email, password, is_admin, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := ur.db.Exec(query, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.IsAdmin, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

// TODO: Validate by email, it returns a boolaean and data i need for login, ask about the fact this function is returning two things, so instead of returning two things, return only a boolean and then create a get endpoint to return all the details of the user
func (ur *userRespository) EmailExists(ctx context.Context, email string) (bool, error) {
	var emailCount int

	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	err := ur.db.QueryRowContext(ctx, query, email).Scan(&emailCount)
	if err != nil {
		fmt.Printf("checking if email exists: %v\n", err)
		return false, err
	}
	return emailCount > 0, nil
}

func (ur *userRespository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := "SELECT * FROM users WHERE email = $1"

	err := ur.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Printf("checking if user exists: %v\n", err)
		return nil, err
	}
	return &user, nil
}
