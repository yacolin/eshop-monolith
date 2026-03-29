package response

import (
	"eshop-monolith/internal/pkg/errcode"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	tid := ""
	if v, ok := c.Get("trace_id"); ok {
		if s, sok := v.(string); sok {
			tid = s
		}
	}
	// ensure trace id is present in response headers as well
	if tid != "" {
		c.Writer.Header().Set("X-Trace-Id", tid)
	}
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: tid,
	})
}

// 业务错误响应
func BizError(c *gin.Context, err *errcode.BizError) {
	status := mapBizErrorToStatus(err)
	tid := ""
	if v, ok := c.Get("trace_id"); ok {
		if s, sok := v.(string); sok {
			tid = s
		}
	}
	c.JSON(status, APIResponse{
		Code:    err.Code,
		Message: err.Message,
		TraceID: tid,
	})
}

// BindError returns a 422 Unprocessable Entity with validation details
func BindError(c *gin.Context, err error) {
	tid := ""
	if v, ok := c.Get("trace_id"); ok {
		if s, sok := v.(string); sok {
			tid = s
		}
	}
	// use business error code for invalid params
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Code:    errcode.ErrInvalidParams.Code,
		Message: "invalid parameters",
		TraceID: tid,
	})
}

// 系统错误响应
func SysError(c *gin.Context, err error) {
	// 记录完整错误到服务端日志，避免把内部错误细节暴露给客户端
	if err != nil {
		log.Printf("sys error: %v, path=%s, method=%s", err, c.Request.URL.Path, c.Request.Method)
	}
	tid := ""
	if v, ok := c.Get("trace_id"); ok {
		if s, sok := v.(string); sok {
			tid = s
		}
	}
	c.JSON(http.StatusInternalServerError, APIResponse{
		Code:    500,
		Message: "internal server error",
		TraceID: tid,
	})
}

func mapBizErrorToStatus(e *errcode.BizError) int {
	switch e {
	case errcode.ErrInvalidParams, errcode.ErrPaginationQuery:
		return http.StatusBadRequest
	case errcode.ErrUnauthorized:
		return http.StatusUnauthorized
	case errcode.ErrProductNotFound, errcode.ErrUserNotFound, errcode.ErrOrderNotFound, errcode.ErrNotFound:
		return http.StatusNotFound
	case errcode.ErrDuplicateOrder:
		return http.StatusConflict
	case errcode.ErrPaymentFailed:
		// payment gateway failure — treat as bad gateway
		return http.StatusBadGateway
	default:
		return http.StatusBadRequest
	}
}
