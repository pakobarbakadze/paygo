package controller

import (
	"net/http"
	"paygo/internal/domain/repository"
	"paygo/internal/domain/service"
	"paygo/internal/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditController struct {
	AuditService *service.AuditService
}

func NewAuditController(db database.DBManager) *AuditController {
	accountRepo := repository.NewAccountRepository(db)
	auditService := service.NewAuditService(accountRepo)

	return &AuditController{
		AuditService: auditService,
	}
}

// AuditAccount godoc
// @Summary Audit a specific account
// @Description Audit an account to detect fraud by comparing ledger entries with transaction history
// @Tags audit
// @Accept json
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} service.AuditResult "Audit result"
// @Failure 400 {object} map[string]interface{} "Invalid account ID"
// @Router /audit/accounts/{accountId} [get]
func (c *AuditController) AuditAccount(ctx *gin.Context) {
	accountIDParam := ctx.Param("accountId")
	accountID, err := uuid.Parse(accountIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	result := c.AuditService.AuditAccount(accountID)

	if result.Status == service.AuditStatusIncomplete {
		ctx.JSON(http.StatusNotFound, result)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
