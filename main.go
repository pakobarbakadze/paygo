package main

import (
	"log"

	"paygo/api/route"
	"paygo/config"
	"paygo/infra/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Setup(&cfg)
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
	defer db.Close()

	r := setupRouter(db)

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(db *database.Database) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	route.SetupRoutes(r, db)

	return r
}
