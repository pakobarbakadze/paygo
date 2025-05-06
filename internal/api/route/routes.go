package route

import (
	"paygo/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *database.Database) {
	v1 := r.Group("/api/v1")

	SetupTransferRoutes(v1, db)
}
