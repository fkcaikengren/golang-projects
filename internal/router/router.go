package router

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-oj/internal/handler"
	"go-oj/internal/pkg/response"
)

func New(healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), requestIDMiddleware())

	r.NoRoute(func(c *gin.Context) {
		resp := response.Error(http.StatusNotFound, "route not found")
		resp.RequestID = c.GetString("request_id")
		c.JSON(http.StatusNotFound, resp)
	})

	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", healthHandler.Health)
		apiV1.GET("/ready", healthHandler.Ready)
	}

	return r
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func newRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)
	}

	return time.Now().UTC().Format("20060102150405.000000000")
}
