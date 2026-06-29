package base

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/user"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type NotificationHandler struct {
	svc *NotificationService
}

func NewNotificationHandler(svc *NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func currentUserID(c *gin.Context) int64 {
	v, _ := c.Get("user_id")
	switch id := v.(type) {
	case int64:
		return id
	case uint:
		return int64(id)
	case float64:
		return int64(id)
	case int:
		return int64(id)
	}
	return 0
}

// ListNotifications 通知列表
// @Summary 通知列表
// @Tags notifications
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=NotificationListResult}
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	page, size := 1, 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(c.Query("size")); err == nil && s > 0 && s <= 100 {
		size = s
	}
	list, total, err := h.svc.ListNotifications(c, currentUserID(c), page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &NotificationListResult{Total: total, List: toRespList(list)})
}

// GetUnreadCount 未读通知数
// @Summary 未读通知数
// @Tags notifications
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=UnreadCountResp}
// @Router /api/v1/notifications/unread [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	count, err := h.svc.GetUnreadCount(c, currentUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &UnreadCountResp{Count: count})
}

// MarkAsRead 标记已读
// @Summary 标记已读
// @Tags notifications
// @Security ApiKeyAuth
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.MarkAsRead(c, id, currentUserID(c)); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// MarkAllAsRead 全部已读
// @Summary 全部已读
// @Tags notifications
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/readall [put]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	if err := h.svc.MarkAllAsRead(c, currentUserID(c)); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// DeleteNotification 删除通知
// @Summary 删除通知
// @Tags notifications
// @Security ApiKeyAuth
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.DeleteNotification(c, id, currentUserID(c)); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// SendSystemNotification 发送系统通知
// @Summary 发送系统通知
// @Tags notifications
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body SendSystemNotificationReq true "系统通知"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/system [post]
func (h *NotificationHandler) SendSystemNotification(c *gin.Context) {
	var req SendSystemNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if _, err := h.svc.CreateNotification(c, req.UserID, req.Title, req.Content, ChannelInApp, CategorySystem); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ── helpers ────────────────────────────────────────

func toResp(n *Notification) *NotificationResp {
	r := &NotificationResp{
		ID:              n.ID,
		UserID:          n.UserID,
		Title:           n.Title,
		Content:         n.Content,
		ContentTemplate: n.ContentTemplate,
		TemplateParams:  n.TemplateParams,
		Channel:         n.Channel,
		Category:        n.Category,
		TargetType:      n.TargetType,
		TargetID:        n.TargetID,
		RedirectURL:     n.RedirectURL,
		IconURL:         n.IconURL,
		IsRead:          n.IsRead,
		IsProcessed:     n.IsProcessed,
		ProcessResult:   n.ProcessResult,
		Priority:        n.Priority,
		CreatedBy:       n.CreatedBy,
		CreatedAt:       n.CreatedAt.UnixMilli(),
		UpdatedAt:       n.UpdatedAt.UnixMilli(),
	}
	if n.ReadAt != nil {
		r.ReadAt = n.ReadAt.UnixMilli()
	}
	if n.ProcessedAt != nil {
		r.ProcessedAt = n.ProcessedAt.UnixMilli()
	}
	return r
}

func toRespList(list []*Notification) []*NotificationResp {
	resp := make([]*NotificationResp, len(list))
	for i, n := range list {
		resp[i] = toResp(n)
	}
	return resp
}

// ── Routes ────────────────────────────────────────

func RegisterNotificationRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) *NotificationService {
	repo := NewNotificationRepository(db)
	svc := NewNotificationService(repo)
	h := NewNotificationHandler(svc)

	notify := v1.Group("/notifications")
	notify.Use(middleware.JWTAuth())
	{
		notify.GET("", h.ListNotifications)
		notify.GET("/unread", h.GetUnreadCount)
		notify.PUT("/:id/read", h.MarkAsRead)
		notify.PUT("/readall", h.MarkAllAsRead)
		notify.DELETE("/:id", h.DeleteNotification)
	}

	roleCfg := user.NewRequireRoleConfig(repos.Role)
	notify.POST("/system", user.RequireAdmin(roleCfg), h.SendSystemNotification)

	return svc
}
