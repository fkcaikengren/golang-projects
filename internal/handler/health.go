package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

type HealthHandler struct {
	service service.HealthService
}

func NewHealthHandler(service service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Health(c *gin.Context) {
	resp := response.Success(h.service.Health(c.Request.Context()))
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if err := h.service.Ready(c.Request.Context()); err != nil {
		resp := response.Error(1001, "database not ready")
		resp.RequestID = c.GetString("request_id")
		c.JSON(http.StatusServiceUnavailable, resp)
		return
	}

	resp := response.Success(map[string]string{"status": "ready"})
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}
