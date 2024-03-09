package main

import (
	"log"

	"github.com/kellswork/wayfarer/internal/api"
	"github.com/kellswork/wayfarer/internal/config"
	"github.com/kellswork/wayfarer/internal/db"
	"github.com/kellswork/wayfarer/internal/db/repositories"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to initiliaze config, %v\n", err.Error())
	}
	// init db
	dbConnection, err := db.ConnectDatabase(cfg.DBURL)
	if err != nil {
		log.Fatalf("unable to initiliaze database, %v\n", err.Error())
	}
	defer db.CloseDatabase(dbConnection)

	repo := repositories.NewRepositories(dbConnection)

	// call server/api
	api.RunServer(repo, cfg)

}
