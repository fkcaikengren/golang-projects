package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

// 定义Handler 依赖service
type AdminAuthHandler struct {
	service service.AdminAuthServiceAPI
}


// DTO
type adminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 创建Handler
func NewAdminAuthHandler(service service.AdminAuthServiceAPI) *AdminAuthHandler {
	return &AdminAuthHandler{service: service}
}


// Handler 挂载方法
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	result, err := h.service.Login(c.Request.Context(), service.AdminLoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	resp := response.Success(result)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *AdminAuthHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeError(c, http.StatusBadRequest, "invalid input")
	case errors.Is(err, service.ErrInvalidCredentials), errors.Is(err, service.ErrInvalidToken):
		writeError(c, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, service.ErrUserDisabled):
		writeError(c, http.StatusForbidden, "user disabled")
	case errors.Is(err, service.ErrInvalidTokenConfig):
		writeError(c, http.StatusInternalServerError, "invalid token config")
	default:
		writeError(c, http.StatusInternalServerError, "internal server error")
	}
}
