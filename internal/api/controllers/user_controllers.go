package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/db/repositories"
	"github.com/kellswork/wayfarer/internal/utils"
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
	Status string             `json:"status"`
	Data   models.CreatedUser `json:"data"`
}

type loginUserResponse struct {
	Status string           `json:"status"`
	Data   models.LoginUser `json:"data"`
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
	doesEmailExists, err := uc.userRepo.EmailExists(c.Request.Context(), userRequestBody.Email)
	if err != nil {
		log.Printf("error occured while verifying if email exist: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{
			Status: "error",
			Error:  "error occured while verifying if email exist",
		})
		return
	}

	if doesEmailExists {
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
		CreatedAt: time.Now(),
	}

	// add user into the database
	if err := uc.userRepo.Create(c.Request.Context(), &user); err != nil {
		// if fail send json failure response
		log.Printf("failed to store user in the db: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{Status: "error", Error: "failed to insert data into the database"})
		return
	}

	// generate json web token and add to header
	token, err := utils.GenerateJwtToken(user.ID)
	if err != nil {
		log.Printf("failed to generate token: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{
			Status: "error",
			Error:  "failed to generate token",
		})
		return
	}
	// set token as header
	c.Header("Authorization", "Bearer "+token)

	// send a response if successful, send json susccess response
	// create a new user and remove the apssword then return that data to the user
	createdUser := models.CreatedUser{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, createUserResponse{Status: "success", Data: createdUser})
}

func (uc *UserControllers) LoginUser(c *gin.Context) {
	// get user input from request body
	var loginUserRequestBody models.LoginUserRequest
	if err := c.BindJSON(&loginUserRequestBody); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())

		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "invalid request body",
		})
		return
	}
	// verify if the email existin teh database
	doesEmailExists, err := uc.userRepo.EmailExists(c.Request.Context(), loginUserRequestBody.Email)
	if err != nil {
		log.Printf("error occured while verifying if email exist: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{
			Status: "error",
			Error:  "error occured while verifying if email exist",
		})
		return
	}

	if !doesEmailExists {
		log.Printf("The email deosnt exists in the db: %v\n", loginUserRequestBody.Email)
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "this email has no account created",
		})
		return
	}
	// decrypt saved hash password and comaper if its the same with the user password
	user, err := uc.userRepo.GetByEmail(c.Request.Context(), loginUserRequestBody.Email)
	if err != nil {
		log.Printf("error occured while getting user details: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, apiResponseError{
			Status: "error",
			Error:  "could not get user details",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserRequestBody.Password))
	if err != nil {
		log.Printf("login password is not correct: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "password is not correct",
		})
		return
	}
	// if the email and password is correct, generate the token
	token, err := utils.GenerateJwtToken(user.ID)
	if err != nil {
		log.Printf("generating jwt token failed: %v\n", loginUserRequestBody.Email)
		c.JSON(http.StatusBadRequest, apiResponseError{
			Status: "error",
			Error:  "failed to generate jwt token",
		})
		return
	}
	// add it to the header
	c.Header("Authorization", "Bearer "+token)
	// return login success message and data accordingly
	loginUser := models.LoginUser{
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
	}
	c.JSON(http.StatusOK, loginUserResponse{Status: "success", Data: loginUser})
}
