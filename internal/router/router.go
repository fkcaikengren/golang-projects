package router

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-oj/internal/handler"
	"go-oj/internal/pkg/response"
)

func New(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	problemSetHandler *handler.ProblemSetHandler,
	problemHandler *handler.ProblemHandler,
	submissionHandler *handler.SubmissionHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(requestIDMiddleware(), gin.Logger(), gin.Recovery())
	r.Use(corsMiddleware())

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

		authV1 := apiV1.Group("/auth")
		{
			authV1.POST("/register", authHandler.Register)
			authV1.POST("/login", authHandler.Login)
		}

		apiV1.GET("/problem-sets", problemSetHandler.List)
		apiV1.GET("/problem-sets/:slug", problemSetHandler.Detail)
		apiV1.GET("/problems", problemHandler.List)
		apiV1.GET("/problems/:slug", problemHandler.Detail)
		apiV1.POST("/submissions", submissionHandler.Submit)
		apiV1.GET("/submissions", submissionHandler.ListMySubmissions)
		apiV1.GET("/problems/:slug/submissions", submissionHandler.ListProblemSubmissions)
	}

	return r
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if !isSafeRequestID(requestID) {
			requestID = newRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func isSafeRequestID(requestID string) bool {
	if requestID == "" || len(requestID) > 64 {
		return false
	}

	for _, char := range requestID {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.", char) {
			continue
		}
		return false
	}

	return true
}

func newRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)
	}

	return time.Now().UTC().Format("20060102150405.000000000")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
