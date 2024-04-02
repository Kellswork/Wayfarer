package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/db/repositories/mocks"
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

}
