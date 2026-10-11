package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/http"
	"adventuria/pkg/ctxlanguage"
	"log"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase/core"

	_ "adventuria/migrations"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	err = ctxlanguage.SetDefaultLang("ru")
	if err != nil {
		log.Fatal(err)
	}

	_, err = adventuria.Start(func(game *adventuria.Game, registry *adventuria.Registry, se *core.ServeEvent) error {
		http.Route(game, registry, se.Router)

		return se.Next()
	})
	if err != nil {
		log.Fatal(err)
	}
}
