package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

type ProblemSetHandler struct {
	service service.ProblemSetServiceAPI
}

func NewProblemSetHandler(service service.ProblemSetServiceAPI) *ProblemSetHandler {
	return &ProblemSetHandler{service: service}
}

func (h *ProblemSetHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := response.Success(items)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *ProblemSetHandler) Detail(c *gin.Context) {
	item, err := h.service.Detail(c.Request.Context(), c.Param("slug"))
	if err != nil {
		if errors.Is(err, service.ErrProblemSetNotFound) {
			writeError(c, http.StatusNotFound, "problem set not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := response.Success(item)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}
