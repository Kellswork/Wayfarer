package repositories

import "database/sql"

type Repositories struct {
	UserRepository UserRepository
	BusRepository  BusRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository: newUserRepository(db),
		BusRepository:  newBusRepository(db),
	}
}
