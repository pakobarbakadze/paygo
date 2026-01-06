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

// AuditAccounts godoc
// @Summary Audit multiple accounts
// @Description Audit multiple accounts concurrently to detect fraud
// @Tags audit
// @Accept json
// @Produce json
// @Param accountIds body []string true "Array of Account IDs"
// @Success 200 {array} service.AuditResult "Array of audit results"
// @Failure 400 {object} map[string]interface{} "Invalid request body or account IDs"
// @Router /audit/accounts [post]
func (c *AuditController) AuditAccounts(ctx *gin.Context) {
	var accountIDStrings []string

	if err := ctx.ShouldBindJSON(&accountIDStrings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if len(accountIDStrings) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one account ID is required"})
		return
	}

	accountIDs := make([]uuid.UUID, 0, len(accountIDStrings))
	for _, idStr := range accountIDStrings {
		accountID, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format", "invalid_id": idStr})
			return
		}
		accountIDs = append(accountIDs, accountID)
	}

	results := c.AuditService.AuditAccounts(accountIDs)

	ctx.JSON(http.StatusOK, gin.H{
		"total":   len(results),
		"results": results,
	})
}
