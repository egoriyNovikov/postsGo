package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Server struct {
	Host string
	Port string
}

func GetDbConfig() DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func GetServerConfig() Server {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return Server{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
}
