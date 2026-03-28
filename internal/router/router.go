package router

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-oj/internal/handler"
	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

func New(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	problemSetHandler *handler.ProblemSetHandler,
	problemHandler *handler.ProblemHandler,
	submissionHandler *handler.SubmissionHandler,
	adminAuthHandler *handler.AdminAuthHandler,
	adminAuthService service.AdminAuthServiceAPI,
	adminHandler *handler.AdminHandler,
	adminAuthorizer *service.AdminAuthorizer,
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
	r.POST("/admin/login", adminAuthHandler.Login)

	admin := r.Group("/admin")
	admin.Use(adminAuthorizationMiddleware(adminAuthService, adminAuthorizer))
	{
		admin.GET("", adminHandler.Dashboard)
		admin.GET("/problems", adminHandler.Problems)
		admin.POST("/problems", adminHandler.Problems)
		admin.PUT("/problems/:id", adminHandler.Problems)
		admin.GET("/problem-sets", adminHandler.ProblemSets)
		admin.POST("/problem-sets", adminHandler.ProblemSets)
		admin.PUT("/problem-sets/:id", adminHandler.ProblemSets)
		admin.GET("/tags", adminHandler.Tags)
		admin.POST("/tags", adminHandler.Tags)
		admin.PUT("/tags/:id", adminHandler.Tags)
		admin.GET("/test-cases", adminHandler.TestCases)
		admin.POST("/test-cases", adminHandler.TestCases)
		admin.PUT("/test-cases/:id", adminHandler.TestCases)
		admin.GET("/judge-configs", adminHandler.JudgeConfigs)
		admin.PUT("/judge-configs/:id", adminHandler.JudgeConfigs)
		admin.GET("/submissions", adminHandler.Submissions)
		admin.GET("/users", adminHandler.Users)
		admin.PATCH("/users/:id/status", adminHandler.Users)
		admin.GET("/settings", adminHandler.Settings)
		admin.PATCH("/settings", adminHandler.Settings)
	}

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
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func adminAuthorizationMiddleware(authService service.AdminAuthServiceAPI, authorizer *service.AdminAuthorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource, action, ok := adminPermissionForRoute(c.Request.Method, c.FullPath())
		if !ok {
			writeAdminError(c, http.StatusForbidden, "forbidden")
			return
		}

		adminUser, err := authService.Authenticate(c.Request.Context(), c.GetHeader("Authorization"))
		if err != nil {
			status := http.StatusUnauthorized
			message := "unauthorized"
			if err == service.ErrUserDisabled {
				status = http.StatusForbidden
				message = "user disabled"
			}
			writeAdminError(c, status, message)
			return
		}

		allowed, err := authorizer.Enforce(strconv.FormatUint(uint64(adminUser.ID), 10), service.AdminDomain, resource, action)
		if err != nil {
			writeAdminError(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if !allowed {
			writeAdminError(c, http.StatusForbidden, "forbidden")
			return
		}

		c.Set("admin_user_id", adminUser.ID)
		c.Next()
	}
}

func adminPermissionForRoute(method string, fullPath string) (string, string, bool) {
	permissions := map[string][2]string{
		"GET /admin":                       {"dashboard", "read"},
		"GET /admin/problems":              {"problems", "read"},
		"POST /admin/problems":             {"problems", "write"},
		"PUT /admin/problems/:id":          {"problems", "write"},
		"GET /admin/problem-sets":          {"problem_sets", "read"},
		"POST /admin/problem-sets":         {"problem_sets", "write"},
		"PUT /admin/problem-sets/:id":      {"problem_sets", "write"},
		"GET /admin/tags":                  {"tags", "read"},
		"POST /admin/tags":                 {"tags", "write"},
		"PUT /admin/tags/:id":              {"tags", "write"},
		"GET /admin/test-cases":            {"test_cases", "read"},
		"POST /admin/test-cases":           {"test_cases", "write"},
		"PUT /admin/test-cases/:id":        {"test_cases", "write"},
		"GET /admin/judge-configs":         {"judge_configs", "read"},
		"PUT /admin/judge-configs/:id":     {"judge_configs", "write"},
		"GET /admin/submissions":           {"submissions", "read"},
		"GET /admin/users":                 {"users", "read"},
		"PATCH /admin/users/:id/status":    {"users", "write"},
		"GET /admin/settings":              {"settings", "read"},
		"PATCH /admin/settings":            {"settings", "write"},
	}

	permission, ok := permissions[method+" "+fullPath]
	return permission[0], permission[1], ok
}

func writeAdminError(c *gin.Context, status int, message string) {
	resp := response.Error(status, message)
	resp.RequestID = c.GetString("request_id")
	c.AbortWithStatusJSON(status, resp)
}
