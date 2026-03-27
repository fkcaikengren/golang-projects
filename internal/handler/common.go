package handler

import (
	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
)

func writeError(c *gin.Context, status int, message string) {
	resp := response.Error(status, message)
	resp.RequestID = c.GetString("request_id")
	c.JSON(status, resp)
}
