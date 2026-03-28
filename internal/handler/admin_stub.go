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

func (h *AdminHandler) Dashboard(c *gin.Context) {
	resp := response.Success(gin.H{
		"title": "Admin Dashboard",
		"admin_user": gin.H{
			"id":            c.GetUint("admin_user_id"),
			"email":         c.GetString("admin_user_email"),
			"display_name":  c.GetString("admin_user_display_name"),
			"status":        c.GetString("admin_user_status"),
			"last_login_at": c.GetInt64("admin_user_last_login_at"),
		},
		"quick_links": []gin.H{
			{"label": "题目管理", "path": "/admin/problems"},
			{"label": "题单管理", "path": "/admin/problem-sets"},
			{"label": "用户管理", "path": "/admin/users"},
			{"label": "系统设置", "path": "/admin/settings"},
		},
	})
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}

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
