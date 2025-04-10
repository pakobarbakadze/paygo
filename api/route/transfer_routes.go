package route

import (
	"paygo/api/controller"

	"github.com/gin-gonic/gin"
)

func SetupTransferRoutes(router *gin.RouterGroup) {
	transferController := controller.NewTransferController()

	transferRoutes := router.Group("/transfers")
	{
		transferRoutes.POST("", transferController.TransferMoney)
	}
}
