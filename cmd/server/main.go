package main

import (
	"log"
	"paygo/internal/api/route"
	"paygo/internal/config"
	"paygo/internal/infra/database"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "paygo/cmd/server/docs"
)

// @title PayGo API
// @version 1.0
// @description Payment and money transfer API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@paygo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
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

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	route.SetupRoutes(r, db)

	return r
}
