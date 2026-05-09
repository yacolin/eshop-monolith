package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// genTraceID 生成唯一的跟踪ID
func genTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("t-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := genTraceID()
		c.Set("trace_id", traceID)
		c.Writer.Header().Set("X-Trace-Id", traceID)
		c.Writer.Header().Set("X-Request-Id", traceID)

		defer func() {
			if rec := recover(); rec != nil {
				err := fmt.Errorf("panic recovered: %v", rec)
				logger.WithRequest(c, "panic recovered",
					"trace_id", traceID,
					"error", err,
					"stack", string(debug.Stack()),
					"method", c.Request.Method,
					"path", c.Request.URL.Path)
				response.SysError(c, err)
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			handleErrors(c, c.Errors.Last().Err, traceID)
		}
	}
}

// handleErrors 处理不同类型的错误
func handleErrors(c *gin.Context, err error, traceID string) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		handleValidationError(c, err, traceID)
		return
	}

	if bizErr, ok := err.(*errcode.BizError); ok {
		handleBusinessError(c, bizErr, traceID)
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		handleBusinessError(c, errcode.ErrNotFound, traceID)
		return
	}

	handleSystemError(c, err, traceID)
}

// handleValidationError 处理验证错误
func handleValidationError(c *gin.Context, err error, traceID string) {
	if gin.Mode() != gin.ReleaseMode {
		logger.WithRequestWarn(c, "validation error",
			"trace_id", traceID,
			"error", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path)
	}
	response.BindError(c, err)
}

// handleBusinessError 处理业务错误
func handleBusinessError(c *gin.Context, bizErr *errcode.BizError, traceID string) {
	logFunc := logger.WithRequest
	logMsg := "business error"
	if bizErr.Code == errcode.ErrUnauthorized.Code {
		logMsg = "authentication error"
	}

	logFunc(c, logMsg,
		"trace_id", traceID,
		"error", bizErr,
		"error_code", bizErr.Code,
		"method", c.Request.Method,
		"path", c.Request.URL.Path)

	// middleware 层负责映射 HTTP 状态，业务层不感知 HTTP
	httpStatus := mapBizErrorToStatus(bizErr)
	response.BizError(c, bizErr, httpStatus)
}

// bizErrorStatusMap 业务错误 → HTTP 状态映射表
// mapBizErrorToStatus 将业务错误码映射为 HTTP 状态码（middleware 层职责）
func mapBizErrorToStatus(e *errcode.BizError) int {
	switch e {
	case errcode.ErrUnauthorized:
		return http.StatusUnauthorized
	case errcode.ErrInsufficientPermissions:
		return http.StatusForbidden
	case errcode.ErrPaymentFailed:
		return http.StatusBadGateway
	case errcode.ErrNotFound, errcode.ErrProductNotFound, errcode.ErrUserNotFound, errcode.ErrOrderNotFound, errcode.ErrPermissionNotFound:
		return http.StatusNotFound
	case errcode.ErrDuplicateOrder, errcode.ErrDuplicateSKU:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

// handleSystemError 处理系统错误
func handleSystemError(c *gin.Context, err error, traceID string) {
	logger.WithRequest(c, "system error",
		"trace_id", traceID,
		"error", err,
		"method", c.Request.Method,
		"path", c.Request.URL.Path)

	response.SysError(c, err)
}
