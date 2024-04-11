package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/db/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {

	// create gin engine
	r := gin.Default()

	// initilaise mocks
	ctrl := gomock.NewController(t)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	// initialise controller with the mock repository
	userctrl := NewUserControllers(mockUserRepo)

	r.POST("/signup", userctrl.CreateUser)

	// create user request body
	sampleUser := models.CreateUserRequest{
		Email:     "kellssyy@gmail.com",
		FirstName: "Kelechi",
		LastName:  "Ogbonna",
		Password:  "1234",
	}

	// check if email exists
	ctx := context.Background()
	mockUserRepo.EXPECT().EmailExists(ctx, sampleUser.Email).Return(false, nil)
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	sampleUserByte, err := json.Marshal(sampleUser)
	require.NoError(t, err)

	require.NoError(t, err)

	// perform HTTP request
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(sampleUserByte))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	// assert the response
	require.Equal(t, http.StatusCreated, recorder.Code)

	result := struct {
		Status string             `json:"status"`
		Data   models.CreatedUser `json:"data"`
	}{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "success", result.Status)
	assert.False(t, result.Data.IsAdmin)

}

func TestLoginUser(t *testing.T) {

	// create gin engine
	r := gin.Default()

	// initilaise mocks
	ctrl := gomock.NewController(t)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	// initialise controller with the mock repository
	userctrl := NewUserControllers(mockUserRepo)

	r.POST("/login", userctrl.LoginUser)

	// create user request body
	sampleUser := models.LoginUserRequest{
		Email:    "kellssyy@gmail.com",
		Password: "1234",
	}

	// check if email exists
	ctx := context.Background()
	mockUserRepo.EXPECT().EmailExists(ctx, sampleUser.Email).Return(true, nil)

	// mock retuned user
	mockReturnedUser := models.User{

		ID:        "e12bab25-e10f-4920-ab57-5fa3dee389ae",
		Email:     sampleUser.Email,
		FirstName: "Kelechi",
		LastName:  "Ogbonna",
		Password:  "$2a$10$th3pXsOR66QBnzRiNlW7We67RfuKlZH2n11JtJ08e4mc7zsMeXXzu", // Hashed password
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}
	mockUserRepo.EXPECT().GetByEmail(ctx, sampleUser.Email).Return(&mockReturnedUser, nil)

	sampleUserByte, err := json.Marshal(sampleUser)
	require.NoError(t, err)

	// perform HTTP request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(sampleUserByte))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	// assert the response
	require.Equal(t, http.StatusOK, recorder.Code)

	result := struct {
		Status string           `json:"status"`
		Data   models.LoginUser `json:"data"`
	}{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "success", result.Status)
	assert.False(t, result.Data.IsAdmin)

}
