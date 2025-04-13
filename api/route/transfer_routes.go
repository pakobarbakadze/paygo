package route

import (
	"paygo/api/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransferRoutes(router *gin.RouterGroup, db *gorm.DB) {
	transferController := controller.NewTransferController(db)

	transferRoutes := router.Group("/transfers")
	{
		transferRoutes.POST("", transferController.TransferMoney)
	}
}
