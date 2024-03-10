package api

import (
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/kellswork/wayfarer/internal/api/controllers"
	"github.com/kellswork/wayfarer/internal/config"
	"github.com/kellswork/wayfarer/internal/db/repositories"
)

// initilaise gin

// create a server struct

// mount your handlers

func RunServer(repo *repositories.Repositories, cfg config.Config) {
	router := gin.Default()

	userController := controllers.NewUserControllers(repo.UserRepository)

	router.POST("/signup", userController.CreateUser)

	// run server in a seperate go routine
	go func() {
		router.Run("localhost:8080")
	}()

	// create channels to listen to shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	// wait for a stop signal to shut down the server
	sig := <-shutdownChan
	log.Printf("shutting down server: %v\n", sig)
}
