package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-oj/internal/pkg/response"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) Dashboard(c *gin.Context)    { h.writeOK(c) }
func (h *AdminHandler) Problems(c *gin.Context)     { h.writeOK(c) }
func (h *AdminHandler) ProblemSets(c *gin.Context)  { h.writeOK(c) }
func (h *AdminHandler) Tags(c *gin.Context)         { h.writeOK(c) }
func (h *AdminHandler) TestCases(c *gin.Context)    { h.writeOK(c) }
func (h *AdminHandler) JudgeConfigs(c *gin.Context) { h.writeOK(c) }
func (h *AdminHandler) Submissions(c *gin.Context)  { h.writeOK(c) }
func (h *AdminHandler) Users(c *gin.Context)        { h.writeOK(c) }
func (h *AdminHandler) Settings(c *gin.Context)     { h.writeOK(c) }

func (h *AdminHandler) writeOK(c *gin.Context) {
	resp := response.Success(gin.H{"ok": true})
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}
