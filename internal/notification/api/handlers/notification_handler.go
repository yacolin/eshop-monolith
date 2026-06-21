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
// @Summary 通知列表
// @Description 分页查询当前用户的通知列表
// @Tags notifications
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=dto.NotificationListResult}
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.NotificationListQuery
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

	response.Success(c, &dto.NotificationListResult{
		Total: result.Total,
		List:  dto.ToNotificationResponseList(result.List),
	})
}

// GetUnreadCount 获取未读通知数
// @Summary 未读通知数
// @Description 获取当前用户的未读通知数量
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.UnreadCountResponse}
// @Router /api/v1/notifications/unread [get]
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

	response.Success(c, &dto.UnreadCountResponse{Count: count})
}

// MarkAsRead 标记单条通知为已读
// @Summary 标记已读
// @Description 标记指定通知为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id}/read [put]
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
// @Summary 全部已读
// @Description 标记当前用户所有通知为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/read-all [put]
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
// @Summary 删除通知
// @Description 删除指定通知
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id} [delete]
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
// @Summary 发送系统通知
// @Description 发送系统通知给指定用户（管理员接口）
// @Tags notifications
// @Accept json
// @Produce json
// @Param body body dto.SendSystemNotificationDTO true "系统通知参数"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/system [post]
func (h *NotificationHandler) SendSystemNotification(c *gin.Context) {
	var req dto.SendSystemNotificationDTO
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
