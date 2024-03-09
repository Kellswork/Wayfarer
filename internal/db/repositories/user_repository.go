package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kellswork/wayfarer/internal/db/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	EmailExists(ctx context.Context, email string) bool
}

type userRespository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) *userRespository {
	return &userRespository{
		db: db,
	}
}

func (ur userRespository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, first_name, last_name, email, password, is_admin) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := ur.db.Exec(query, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.IsAdmin)

	if err != nil {
		return err
	}
	return nil
}

// TODO: Validate by email
func (ur userRespository) EmailExists(ctx context.Context, email string) bool {
	var emailCount int
	query := "SELECT * FROM users WHERE email = $1"

	err := ur.db.QueryRow(query, email).Scan(&emailCount)
	if err != nil {
		fmt.Println("checking if email exists:", err)
		return false
	}
	return emailCount > 0
}
