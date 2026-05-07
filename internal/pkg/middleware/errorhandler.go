package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"eshop-monolith/internal/pkg/errcode"
	"eshop-monolith/internal/pkg/logger"
	"eshop-monolith/internal/pkg/response"

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
		// 生成跟踪ID并设置到上下文和响应头
		traceID := genTraceID()
		c.Set("trace_id", traceID)
		c.Writer.Header().Set("X-Trace-Id", traceID)
		c.Writer.Header().Set("X-Request-Id", traceID) // 增加Request-ID头

		// 捕获panic
		defer func() {
			if rec := recover(); rec != nil {
				// 构造错误并记录详细日志
				err := fmt.Errorf("panic recovered: %v", rec)
				logger.WithRequest(c, "panic recovered",
					"trace_id", traceID,
					"error", err,
					"stack", string(debug.Stack()),
					"method", c.Request.Method,
					"path", c.Request.URL.Path)
				// 返回系统错误响应
				response.SysError(c, err)
				// 确保请求被中止
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()

		// 处理记录的错误
		if len(c.Errors) > 0 {
			handleErrors(c, c.Errors.Last().Err, traceID)
		}
	}
}

// handleErrors 处理不同类型的错误
func handleErrors(c *gin.Context, err error, traceID string) {
	// 验证错误
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		handleValidationError(c, err, traceID)
		return
	}

	// 业务错误
	if bizErr, ok := err.(*errcode.BizError); ok {
		handleBusinessError(c, bizErr, traceID)
		return
	}

	// GORM 记录未找到 → 404 业务错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handleBusinessError(c, errcode.ErrNotFound, traceID)
		return
	}

	// 系统错误
	handleSystemError(c, err, traceID)
}

// handleValidationError 处理验证错误
func handleValidationError(c *gin.Context, err error, traceID string) {
	// 在开发环境中记录验证错误
	if gin.Mode() != gin.ReleaseMode {
		logger.WithRequestWarn(c, "validation error",
			"trace_id", traceID,
			"error", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path)
	}
	// 返回422错误
	response.BindError(c, err)
}

// handleBusinessError 处理业务错误
func handleBusinessError(c *gin.Context, bizErr *errcode.BizError, traceID string) {
	// 根据错误类型记录不同级别的日志
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

	// 返回业务错误响应
	response.BizError(c, bizErr)
}

// handleSystemError 处理系统错误
func handleSystemError(c *gin.Context, err error, traceID string) {
	// 记录详细的系统错误日志
	logger.WithRequest(c, "system error",
		"trace_id", traceID,
		"error", err,
		"method", c.Request.Method,
		"path", c.Request.URL.Path)

	// 返回500错误
	response.SysError(c, err)
}
