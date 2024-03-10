package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/db/repositories"
	"golang.org/x/crypto/bcrypt"
)

// conventionally, we create types then create models under
type UserControllers struct {
	userRepo repositories.UserRepository
}

type apiResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type createUserResponse struct {
	Status string      `json:"status"`
	Data   models.User `json:"data"`
}

func NewUserControllers(userRepo repositories.UserRepository) *UserControllers {
	return &UserControllers{
		userRepo: userRepo,
	}
}

func (uc *UserControllers) CreateUser(c *gin.Context) {
	// validate request body
	var userRequestBody models.CreateUserRequest
	if err := c.BindJSON(&userRequestBody); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())

		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "invalid request body",
		})

		return
	}
	// validate the user input
	validate := validator.New()
	err := validate.Struct(userRequestBody)
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
	result := uc.userRepo.EmailExists(c.Request.Context(), userRequestBody.Email)
	if result {
		log.Printf("The email already exists: %v\n", userRequestBody.Email)

		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "This email already exists",
		})
		return
	}

	// hash passowrd
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequestBody.Password), 10)
	if err != nil {
		log.Printf("failed to hash password: %v\n", err.Error())

		c.JSON(http.StatusInternalServerError, apiResponseError{
			Status: "error",
			Error:  "could not handle request",
		})
		return
	}
	// add details to sql table
	user := models.User{
		ID:        uuid.NewString(),
		Email:     userRequestBody.Email,
		FirstName: userRequestBody.FirstName,
		LastName:  userRequestBody.LastName,
		Password:  string(hashedPassword),
		IsAdmin:   false,
	}

	// add user into the database
	if err := uc.userRepo.Create(c.Request.Context(), &user); err != nil {
		// if fail send json failure response
		log.Printf("failed to store user in the db: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to insert data into the database"})
		return
	}
	// send a response if successful, send json susccess response
	c.JSON(http.StatusCreated, createUserResponse{Status: "success", Data: user})
}
