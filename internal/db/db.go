package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // undscore to indicate it is noit really used but has to be there
)

// setup postgres databse connection
// the function will either return an sql.db driver or an error
func ConnectDatabase(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("an error occured while trying to connect to the database")
		return nil, err // return nil for the *sql.DB and the error
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("database connected")
	return db, nil // return the db connection and nil for the error

}

// call it with defer at the end before fun main return function, to close the database when the program exists
func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing the database connection: %v", err)
	}
}
