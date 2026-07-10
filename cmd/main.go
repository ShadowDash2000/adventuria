package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/http"
	"log"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase/core"

	_ "adventuria/migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	_, err := adventuria.Start(func(game *adventuria.Game, se *core.ServeEvent) error {
		http.Route(game, se.Router)

		return se.Next()
	})
	if err != nil {
		log.Fatal(err)
	}
}
