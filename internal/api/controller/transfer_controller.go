package controller

import (
	"net/http"
	"paygo/internal/api/dto"
	"paygo/internal/domain/repository"
	"paygo/internal/domain/service"
	"paygo/internal/infra/database"

	"github.com/gin-gonic/gin"
)

type TransferController struct {
	TransferService *service.TransferService
}

func NewTransferController(db database.DBManager) *TransferController {
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	transferService := service.NewTransferService(db, accountRepo, transactionRepo)

	return &TransferController{
		TransferService: transferService,
	}
}

// TransferMoney godoc
// @Summary Transfer money between accounts
// @Description Transfer money from one account to another
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body dto.TransferRequest true "Transfer details"
// @Success 200 {object} map[string]interface{} "Transfer successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /transfers [post]
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
