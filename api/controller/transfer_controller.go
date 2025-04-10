package controller

import (
	"net/http"
	"paygo/api/dto"
	"paygo/domain/service"
	"paygo/infra/database"

	"github.com/gin-gonic/gin"
)

type TransferController struct {
	TransferService *service.TransferService
}

func NewTransferController() *TransferController {
	transferService := service.NewTransferService(database.DB)

	return &TransferController{
		TransferService: transferService,
	}
}

func (c *TransferController) TransferMoney(ctx *gin.Context) {
	var request dto.TransferRequest

	println("In transfer money controller")

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.TransferService.TransferMoney(
		request.FromAccountID,
		request.ToAccountID,
		request.Amount,
		request.Description,
	)
}
