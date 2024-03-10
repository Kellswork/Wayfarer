package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// struct to represent json webtoken

func generateECDSAPrivateKey() (*ecdsa.PrivateKey, error) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("failed to generate private key: %v", err.Error())
		return nil, err
	}
	return privateKey, nil
}

func GenerateJwtToken(ID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":  ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtKey, _ := generateECDSAPrivateKey()
	return token.SignedString(jwtKey)
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	jwtKey, _ := generateECDSAPrivateKey()
	token, err := jwt.Parse(tokenString, func(jt *jwt.Token) (interface{}, error) {

		_, ok := jt.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signin method %v", jt.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}
