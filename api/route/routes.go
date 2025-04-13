package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	SetupTransferRoutes(v1, db)
}
