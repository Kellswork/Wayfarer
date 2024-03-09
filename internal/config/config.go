package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort   = "8080"
	defaultDburl  = "postgresql://localhost:5432"
	defaultAppEnv = "local"
)

type Config struct {
	DBURL  string
	Port   string
	DBName string
	AppEnv string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, err
	}

	config := Config{}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}
	config.Port = port

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok {
		return Config{}, fmt.Errorf("error no db url")
	}
	config.DBURL = dbURL

	appEnv, ok := os.LookupEnv("APP_ENV")
	if !ok {
		appEnv = defaultAppEnv
	}
	config.AppEnv = appEnv

	return config, nil
}

// connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
