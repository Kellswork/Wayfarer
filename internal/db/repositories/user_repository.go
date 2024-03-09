package repositories

import (
	"database/sql"

	"github.com/kellswork/wayfarer/internal/db/models"
)

type UserRepository interface {
	Create(user *models.User) error
}

type userRespository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRespository {
	return &userRespository{
		db: db,
	}
}

func (ur userRespository) Create(user *models.User) error {
	query := "INSERT INTO users (first_name, last_name, email, password, is_admin) VALUES ($1, $2, $3, $4, $5)"

	_, err := ur.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Email, user.Password, user.IsAdmin)

	if err != nil {
		return err
	}
	return nil
}
