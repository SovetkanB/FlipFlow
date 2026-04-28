package main

import (
	"log"

	"github.com/SovetkanB/FlipFlow/config"
	"github.com/SovetkanB/FlipFlow/internal/repo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using env variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	db, err := repo.NewDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")
}
