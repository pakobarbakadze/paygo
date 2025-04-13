package controller

import (
	"net/http"
	"paygo/api/dto"
	"paygo/domain/repository"
	"paygo/domain/service"
	"paygo/infra/database"

	"github.com/gin-gonic/gin"
)

type TransferController struct {
	TransferService *service.TransferService
}

func NewTransferController(db *database.Database) *TransferController {
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	transferService := service.NewTransferService(db, accountRepo, transactionRepo)

	return &TransferController{
		TransferService: transferService,
	}
}

func (c *TransferController) TransferMoney(ctx *gin.Context) {
	var request dto.TransferRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, fromAccount, toAccount, err := c.TransferService.TransferMoney(
		request.FromAccountID,
		request.ToAccountID,
		request.Amount,
		request.Description,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dto.TransferResponse{
		TransactionID:         transaction.ID,
		TransactionReference:  transaction.TransactionReference,
		Status:                transaction.Status,
		Amount:                transaction.Amount,
		CurrencyCode:          transaction.CurrencyCode,
		FromAccountID:         fromAccount.ID,
		ToAccountID:           toAccount.ID,
		FromAccountNewBalance: fromAccount.Balance,
		ToAccountNewBalance:   toAccount.Balance,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
		"data":    response,
	})
}
