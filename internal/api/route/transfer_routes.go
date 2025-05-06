package route

import (
	"paygo/internal/api/controller"
	"paygo/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func SetupTransferRoutes(router *gin.RouterGroup, db *database.Database) {
	transferController := controller.NewTransferController(db)

	transferRoutes := router.Group("/transfers")
	{
		transferRoutes.POST("", transferController.TransferMoney)
	}
}
