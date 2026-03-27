package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

type AuthHandler struct {
	service service.AuthServiceAPI
}

func NewAuthHandler(service service.AuthServiceAPI) *AuthHandler {
	return &AuthHandler{service: service}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	result, err := h.service.Register(c.Request.Context(), service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	resp := response.Success(result)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	result, err := h.service.Login(c.Request.Context(), service.LoginInput{
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

func (h *AuthHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeError(c, http.StatusBadRequest, "invalid input")
	case errors.Is(err, service.ErrDuplicateEmail):
		writeError(c, http.StatusConflict, "duplicate email")
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(c, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, service.ErrUserDisabled):
		writeError(c, http.StatusForbidden, "user disabled")
	case errors.Is(err, service.ErrInvalidTokenConfig):
		writeError(c, http.StatusInternalServerError, "invalid token config")
	default:
		writeError(c, http.StatusInternalServerError, "internal server error")
	}
}
