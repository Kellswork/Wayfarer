package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {

	// set go gin as router
	router := gin.Default()

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
