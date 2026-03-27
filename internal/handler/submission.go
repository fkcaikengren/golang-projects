package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

type SubmissionHandler struct {
	service service.SubmissionServiceAPI
}

type submitRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	ProblemID uint   `json:"problem_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

func NewSubmissionHandler(service service.SubmissionServiceAPI) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

func (h *SubmissionHandler) Submit(c *gin.Context) {
	var req submitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	item, err := h.service.Submit(c.Request.Context(), service.SubmitCodeInput{
		UserID:    req.UserID,
		ProblemID: req.ProblemID,
		Language:  req.Language,
		Code:      req.Code,
	})
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	resp := response.Success(item)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *SubmissionHandler) ListMySubmissions(c *gin.Context) {
	userID, err := parseUintQuery(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	items, err := h.service.ListUserSubmissions(c.Request.Context(), userID)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	resp := response.Success(items)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *SubmissionHandler) ListProblemSubmissions(c *gin.Context) {
	userID, err := parseUintQuery(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	items, err := h.service.ListProblemSubmissions(c.Request.Context(), userID, c.Param("slug"))
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	resp := response.Success(items)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *SubmissionHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeError(c, http.StatusBadRequest, "invalid input")
	case errors.Is(err, service.ErrProblemNotFound):
		writeError(c, http.StatusNotFound, "problem not found")
	default:
		writeError(c, http.StatusInternalServerError, "internal server error")
	}
}

func parseUintQuery(value string) (uint, error) {
	if value == "" {
		return 0, errors.New("missing")
	}

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid")
	}
	return uint(id), nil
}
