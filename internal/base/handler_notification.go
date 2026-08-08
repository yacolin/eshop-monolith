package base

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/repository"
	ws "eshop-monolith/internal/infra/ws"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type NotificationHandler struct {
	svc *NotificationService
	hub *ws.Hub
}

func NewNotificationHandler(svc *NotificationService, h *ws.Hub) *NotificationHandler {
	return &NotificationHandler{svc: svc, hub: h}
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
	var req NotificationListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListNotifications(c, currentUserID(c), &req)
	if err != nil {
		c.Error(err)
		return
	}
	// response.Success(c, &NotificationListResult{Total: total, List: toRespList(list)})
	response.Success(c, result)
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
	n, err := h.svc.CreateSystemNotification(c, req.UserID, req.TemplateCode, req.Title, req.Content)
	if err != nil {
		c.Error(err)
		return
	}
	// WebSocket 实时推送
	if h.hub != nil {
		if req.UserID == 0 {
			h.broadcastToAll(n)
		} else {
			h.pushToWS(req.UserID, n)
		}
	}
	response.Success(c, nil)
}

// ListTemplates 通知模板列表
// @Summary 通知模板列表
// @Tags notifications
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=NotificationTemplateListResult}
// @Router /api/v1/notifications/templates [get]
func (h *NotificationHandler) ListTemplates(c *gin.Context) {
	list, err := h.svc.ListTemplates(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &NotificationTemplateListResult{
		List: templateToRespList(list),
	})
}

// pushToWS 通过 WebSocket 推送通知给指定用户
func (h *NotificationHandler) pushToWS(userID int64, n *Notification) {
	msg := h.buildWSMessage(n)
	data, _ := json.Marshal(msg)
	h.hub.SendToUser(userID, data)
}

// broadcastToAll 通过 WebSocket 广播通知给所有在线用户
func (h *NotificationHandler) broadcastToAll(n *Notification) {
	msg := h.buildWSMessage(n)
	data, _ := json.Marshal(msg)
	h.hub.Broadcast(data)
}

func (h *NotificationHandler) buildWSMessage(n *Notification) map[string]interface{} {
	return map[string]interface{}{
		"seq":   h.hub.NextSeq(),
		"type":  "notification",
		"payload": map[string]interface{}{
			"id":       n.ID,
			"title":    n.Title,
			"content":  n.Content,
			"channel":  n.Channel,
			"category": n.Category,
			"is_read":  false,
		},
	}
}

// ── helpers ────────────────────────────────────────

func toResp(n *Notification) *NotificationResp {
	return &NotificationResp{
		ID:              n.ID,
		UserID:          n.UserID,
		Title:           n.Title,
		Content:         n.Content,
		ContentTemplate: n.ContentTemplate,
		TemplateParams:  safeStr(n.TemplateParams),
		Channel:         n.Channel,
		Category:        n.Category,
		TargetType:      n.TargetType,
		TargetID:        n.TargetID,
		RedirectURL:     n.RedirectURL,
		IconURL:         n.IconURL,
		IsRead:          n.IsRead,
		Priority:        n.Priority,
		CreatedBy:       n.CreatedBy,
		CreatedAt:       n.CreatedAt.UnixMilli(),
		UpdatedAt:       n.UpdatedAt.UnixMilli(),
	}
}

func toRespList(list []*Notification) []*NotificationResp {
	resp := make([]*NotificationResp, len(list))
	for i, n := range list {
		resp[i] = toResp(n)
	}
	return resp
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func templateToResp(t *NotificationTemplate) *NotificationTemplateResp {
	return &NotificationTemplateResp{
		ID:              t.ID,
		TemplateCode:    t.TemplateCode,
		Channel:         t.Channel,
		TitleTemplate:   t.TitleTemplate,
		ContentTemplate: t.ContentTemplate,
		Category:        t.Category,
		Priority:        t.Priority,
		Status:          t.Status,
	}
}

func templateToRespList(list []NotificationTemplate) []*NotificationTemplateResp {
	resp := make([]*NotificationTemplateResp, len(list))
	for i, t := range list {
		resp[i] = templateToResp(&t)
	}
	return resp
}

// ── Routes ────────────────────────────────────────

func RegisterNotificationRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, wsHub *ws.Hub) *NotificationService {
	repo := NewNotificationRepository(db)
	svc := NewNotificationService(repo)
	h := NewNotificationHandler(svc, wsHub)


	notify := v1.Group("/notifications")
	notify.Use(middleware.JWTAuth())
	{
		notify.GET("", h.ListNotifications)
		notify.GET("/templates", h.ListTemplates)
		notify.GET("/unread", h.GetUnreadCount)
		notify.PUT("/:id/read", h.MarkAsRead)
		notify.PUT("/readall", h.MarkAllAsRead)
		notify.DELETE("/:id", h.DeleteNotification)
	}

	notify.POST("/system", middleware.RequireAdmin(), h.SendSystemNotification)

	return svc
}
