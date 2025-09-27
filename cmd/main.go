package main

import (
	"fmt"
	"log"

	"github.com/egoriyNovikov/internal/config"
	"github.com/egoriyNovikov/internal/db"
	"github.com/egoriyNovikov/internal/router"
)

func main() {
	fmt.Println("Hello, World!")
	db, err := db.Connect(config.GetDbConfig())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	fmt.Println("Connected to database")

	router.Router(db, config.GetServerConfig().Port)

	log.Println("Server started on port " + config.GetServerConfig().Port)
}
