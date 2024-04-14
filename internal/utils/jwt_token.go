package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtPrivateKey []byte

func getJwtToken() {

	var err = godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load env file: %v\n", err.Error())
	}
	jwtPrivateKey = []byte(os.Getenv("JWT_SECRET"))
}

func GenerateJwtToken(ID string, IsAdmin bool) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      ID,
		"isAdmin": IsAdmin,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	getJwtToken()
	ss, err := token.SignedString(jwtPrivateKey)
	return ss, err
}

func VerifyJwtToken(tokenString string) (jwt.MapClaims, error) {
	getJwtToken()
	token, err := jwt.Parse(tokenString, func(jt *jwt.Token) (interface{}, error) {

		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", jt.Header["alg"])
		}

		return jwtPrivateKey, nil

	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
