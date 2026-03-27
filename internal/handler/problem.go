package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
	"go-oj/internal/service"
)

type ProblemHandler struct {
	service service.ProblemServiceAPI
}

type problemListQuery struct {
	Difficulty string `form:"difficulty"`
	Tag        string `form:"tag"`
	Keyword    string `form:"keyword"`
}

func NewProblemHandler(service service.ProblemServiceAPI) *ProblemHandler {
	return &ProblemHandler{service: service}
}

func (h *ProblemHandler) List(c *gin.Context) {
	var query problemListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input")
		return
	}

	items, err := h.service.List(c.Request.Context(), service.ProblemListInput{
		Difficulty: query.Difficulty,
		Tag:        query.Tag,
		Keyword:    query.Keyword,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := response.Success(items)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

func (h *ProblemHandler) Detail(c *gin.Context) {
	item, err := h.service.Detail(c.Request.Context(), c.Param("slug"))
	if err != nil {
		if errors.Is(err, service.ErrProblemNotFound) {
			writeError(c, http.StatusNotFound, "problem not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := response.Success(item)
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}
