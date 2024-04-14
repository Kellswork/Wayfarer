package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/db/repositories"
)

// conventionally, we create types then create models under
type TripControllers struct {
	tripRepo repositories.TripRepository
}

type apiTripResponseSucess struct {
	Status string      `json:"status"`
	Data   models.Trip `json:"data"`
}

type apiGetTripsResponseSucess struct {
	Status string        `json:"status"`
	Data   []models.Trip `json:"data"`
}

func NewTripControllers(tripRepo repositories.TripRepository) *TripControllers {
	return &TripControllers{
		tripRepo: tripRepo,
	}
}

func (tc *TripControllers) CreateTrip(c *gin.Context) {
	// validate request body
	var tripReqBody models.TripReqBody

	if err := c.BindJSON(&tripReqBody); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())

		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "invalid request body",
		})

		return
	}

	validate := validator.New()
	err := validate.Struct(tripReqBody)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Printf("request body valiodation failed: %v\n", validationErrors.Error())
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "Invalid request body: " + validationErrors.Error(),
		})
		return
	}

	// add details to sql table
	trip := models.Trip{
		ID:          uuid.NewString(),
		BusID:       tripReqBody.BusID,
		Origin:      tripReqBody.Origin,
		Destination: tripReqBody.Destination,
		TripDate:    tripReqBody.TripDate,
		Fare:        tripReqBody.Fare,
		Status:      tripReqBody.Status,
		CreatedAt:   time.Now(),
	}

	// add user into the database
	err = tc.tripRepo.Create(c.Request.Context(), &trip)

	if err != nil {
		// if fail send json failure response
		log.Printf("failed to add trip to the db: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to insert trip data into the database"})
		return
	}

	c.JSON(http.StatusOK, apiTripResponseSucess{Status: "success", Data: trip})
}

func (tc *TripControllers) CancelTrip(c *gin.Context) {

	tripID := c.Param("tripID")
	fmt.Println("Trip ID:", tripID)

	var tripReqBody models.CancelTripReqBody

	if err := c.BindJSON(&tripReqBody); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())

		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "invalid request body",
		})

		return
	}

	validate := validator.New()
	err := validate.Struct(tripReqBody)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Printf("request body valiodation failed: %v\n", validationErrors.Error())
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "Invalid request body: " + validationErrors.Error(),
		})
		return
	}

	trip := models.CancelTrip{
		Status:    tripReqBody.Status,
		UpdatedAt: time.Now(),
	}

	err = tc.tripRepo.Cancel(c.Request.Context(), string(trip.Status), trip.UpdatedAt, tripID)

	if err != nil {
		log.Printf("failed to cancel trip: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to update trip status in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "trip status updated to 'canceled'"})

}

func (tc *TripControllers) GetAllTrips(c *gin.Context) {

	// add user into the database
	result, err := tc.tripRepo.FindAll(c.Request.Context())

	if err != nil {
		// if fail send json failure response
		log.Printf("failed to get trips from the db: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to get trips from the db"})
		return
	}

	c.JSON(http.StatusOK, apiGetTripsResponseSucess{Status: "sucess", Data: *result})
}
