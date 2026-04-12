package main

import (
	"log"

	"practice-7/config"
	"practice-7/internal/app"
)

func main() {
	cfg := config.Load()
	if err := app.Run(cfg); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
