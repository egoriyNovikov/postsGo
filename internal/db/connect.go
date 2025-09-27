package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/egoriyNovikov/internal/config"
	_ "github.com/lib/pq"
)

func Connect(db config.DB) (*sql.DB, error) {
	connect, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.DBName))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	}
	return connect, nil
}
