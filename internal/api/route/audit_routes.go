package route

import (
	"paygo/internal/api/controller"
	"paygo/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func SetupAuditRoutes(rg *gin.RouterGroup, db database.DBManager) {
	auditController := controller.NewAuditController(db)

	auditGroup := rg.Group("/audit")
	{
		auditGroup.GET("/accounts/:accountId", auditController.AuditAccount)
		auditGroup.POST("/accounts", auditController.AuditAccounts)
	}
}
