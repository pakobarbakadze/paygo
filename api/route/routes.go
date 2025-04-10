package route

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	SetupTransferRoutes(v1)
}
