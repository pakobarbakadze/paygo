package route

import (
	"paygo/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func SetupHealthRoutes(router *gin.RouterGroup) {
	healthController := controller.NewHealthController()

	router.GET("/ping", healthController.Ping)
}
