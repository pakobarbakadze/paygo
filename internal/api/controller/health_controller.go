package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// Ping godoc
// @Summary Health check endpoint
// @Description Returns pong to verify the server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (h *HealthController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
