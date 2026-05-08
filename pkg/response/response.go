package response

import (
	"errors"
	"log"
	"net/http"

	"eshop-monolith/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// ValidationFieldError 字段验证错误详情
type ValidationFieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
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

// BindError returns a 422 Unprocessable Entity with structured field validation details
func BindError(c *gin.Context, err error) {
	tid := ""
	if v, ok := c.Get("trace_id"); ok {
		if s, sok := v.(string); sok {
			tid = s
		}
	}

	// 尝试提取字段级验证错误详情
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make([]ValidationFieldError, len(ve))
		for i, fe := range ve {
			details[i] = ValidationFieldError{
				Field:   fe.Field(),
				Tag:     fe.Tag(),
				Message: buildValidationMessage(fe),
			}
		}
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Code:    errcode.ErrInvalidParams.Code,
			Message: "invalid parameters",
			Data:    details,
			TraceID: tid,
		})
		return
	}

	// 非验证错误，保持原有逻辑
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Code:    errcode.ErrInvalidParams.Code,
		Message: "invalid parameters",
		TraceID: tid,
	})
}

// buildValidationMessage 根据验证规则生成友好的错误消息
func buildValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "'" + fe.Field() + "' is required"
	case "max":
		return "'" + fe.Field() + "' must be at most " + fe.Param()
	case "min":
		return "'" + fe.Field() + "' must be at least " + fe.Param()
	case "gt":
		return "'" + fe.Field() + "' must be greater than " + fe.Param()
	case "gte":
		return "'" + fe.Field() + "' must be at least " + fe.Param()
	case "lte":
		return "'" + fe.Field() + "' must be at most " + fe.Param()
	case "email":
		return "'" + fe.Field() + "' must be a valid email address"
	case "oneof":
		return "'" + fe.Field() + "' must be one of: " + fe.Param()
	case "len":
		return "'" + fe.Field() + "' must be exactly " + fe.Param() + " characters"
	case "dive":
		return "'" + fe.Field() + "' contains an invalid item"
	default:
		return fe.Error()
	}
}

// 系统错误响应
func SysError(c *gin.Context, err error) {
	// 记录完整错误到服务端日志
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
		Message: err.Error(),
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
	case errcode.ErrDuplicateSKU:
		return http.StatusConflict
	case errcode.ErrPaymentFailed:
		// payment gateway failure — treat as bad gateway
		return http.StatusBadGateway
	default:
		return http.StatusBadRequest
	}
}
