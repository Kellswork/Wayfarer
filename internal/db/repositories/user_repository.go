package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kellswork/wayfarer/internal/db/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	EmailExists(ctx context.Context, email string) (models.User, bool)
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

// TODO: Validate by email, it returns a boolaean and data i need for login, ask about the fact this function is returning two things, so instead of returning two things, return only a boolean and then create a get endpoint to return all the details of the user
func (ur userRespository) EmailExists(ctx context.Context, email string) (models.User, bool) {
	var user models.User

	query := "SELECT id, is_admin, email, password FROM users WHERE email = $1"

	err := ur.db.QueryRow(query, email).Scan(&user.ID, &user.IsAdmin, &user.Email, &user.Password)
	if err != nil {
		fmt.Printf("checking if email exists: %v\n", err)
		return user, false
	}
	return user, len(user.Email) > 0
}
