package main

import (
	"log"
	"paygo/internal/config"
	"paygo/internal/infra/database"
)

func main() {
	log.Println("PayGo Database Seeder")
	log.Println("=====================")

	cfg := config.LoadConfig()

	db, err := database.Setup(&cfg)
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
	defer db.Close()

	if err := db.Seed(); err != nil {
		log.Fatalf("Database seeding failed: %v", err)
	}

	log.Println("\n✓ Database seeded successfully!")
}
