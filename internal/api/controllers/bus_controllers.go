package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/db/repositories"
)

// conventionally, we create types then create models under
type BusControllers struct {
	busRepo repositories.BusRepository
}

type apiBusResponseSucess struct {
	Status string             `json:"status"`
	Data   *models.CreatedBus `json:"data"`
}

type getAllBusesResponse struct {
	Status string              `json:"status"`
	Data   []models.CreatedBus `json:"data"`
}

func NewBusControllers(busRepo repositories.BusRepository) *BusControllers {
	return &BusControllers{
		busRepo: busRepo,
	}
}

func (bc *BusControllers) AddBus(c *gin.Context) {
	// validate request body
	var busReqBody models.BusReqBody

	if err := c.BindJSON(&busReqBody); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())

		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "invalid request body",
		})

		return
	}
	// validate the user input
	validate := validator.New()
	err := validate.Struct(busReqBody)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Printf("request body valiodation failed: %v\n", validationErrors.Error())
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "Invalid request body: " + validationErrors.Error(),
		})
		return
	}

	// check if the email has been used before and return a json response accordingly
	doesPlateNumberExist, err := bc.busRepo.DoesPlateExists(c.Request.Context(), busReqBody.PlateNumber)
	if err != nil {
		log.Printf("error occured while verifying if plate number exist: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{
			Status: "error",
			Error:  "error occured while verifying if plate number exist",
		})
		return
	}

	if doesPlateNumberExist {
		log.Printf("This plate number already exists: %v\n", busReqBody.PlateNumber)
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "This plate number already exists",
		})
		return
	}

	// add details to sql table
	bus := models.Bus{
		PlateNumber:  busReqBody.PlateNumber,
		Manufacturer: busReqBody.Manufacturer,
		Model:        busReqBody.Model,
		Type:         busReqBody.Type,
		Year:         busReqBody.Year,
		Capacity:     busReqBody.Capacity,
		CreatedAt:    time.Now(),
	}

	// add user into the database
	result, err := bc.busRepo.Create(c.Request.Context(), &bus)

	if err != nil {
		// if fail send json failure response
		log.Printf("failed to add bus to the db: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to insert bus data into the database"})
		return
	}
	fmt.Println()

	c.JSON(http.StatusOK, apiBusResponseSucess{Status: "success", Data: result})
}

func (bc *BusControllers) GetAllBuses(c *gin.Context) {

	// add user into the database
	result, err := bc.busRepo.FindAllBuses(c.Request.Context())

	if err != nil {
		// if fail send json failure response
		log.Printf("failed to fetch buses from the db: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to fetch buses from the database"})
		return
	}

	c.JSON(http.StatusOK, getAllBusesResponse{Status: "sucess", Data: *result})
}
