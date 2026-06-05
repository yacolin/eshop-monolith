package handlers

import (
	"strconv"

	"eshop-monolith/internal/notification/api/dto"
	"eshop-monolith/internal/notification/domain/models"
	"eshop-monolith/internal/notification/service"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	svc *service.NotificationService
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// ListNotifications 获取通知列表（分页）
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.ListNotificationReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.Pagination.Normalize()

	result, err := h.svc.ListNotifications(c, userID, req.Page, req.Size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, &dto.NotificationListResp{
		Total: result.Total,
		List:  dto.ToNotificationRespList(result.List),
	})
}

// GetUnreadCount 获取未读通知数
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	count, err := h.svc.GetUnreadCount(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, &dto.UnreadCountResp{Count: count})
}

// MarkAsRead 标记单条通知为已读
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	notificationID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.MarkAsRead(c, notificationID, userID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "已标记为已读"})
}

// MarkAllAsRead 标记所有通知为已读
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.MarkAllAsRead(c, userID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "全部标记为已读"})
}

// DeleteNotification 删除通知
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	notificationID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.DeleteNotification(c, notificationID, userID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "已删除"})
}

// SendSystemNotification 发送系统通知（管理员接口）
func (h *NotificationHandler) SendSystemNotification(c *gin.Context) {
	var req dto.SendSystemNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	_, err := h.svc.CreateNotification(c, req.UserID, req.Title, req.Content, models.NotificationTypeSystem)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "系统通知发送成功"})
}

// getCurrentUserID 从 Gin Context 获取当前用户 ID
func getCurrentUserID(c *gin.Context) (int64, error) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, errcode.ErrUnauthorized
	}
	switch id := v.(type) {
	case float64:
		return int64(id), nil
	case uint:
		return int64(id), nil
	case int64:
		return id, nil
	case int:
		return int64(id), nil
	case string:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return 0, errcode.ErrUnauthorized
		}
		return n, nil
	default:
		return 0, errcode.ErrUnauthorized
	}
}
