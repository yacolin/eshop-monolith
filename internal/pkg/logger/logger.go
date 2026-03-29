package logger

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var SugaredLogger *zap.SugaredLogger

func init() {
	SugaredLogger = NewProductionLogger()
}

func NewProductionLogger() *zap.SugaredLogger {
	// 创建日志目录（如果不存在）
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// 设置编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "production" // 默认为生产环境
	}

	// 设置编码器
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台输出
	consoleSyncer := zapcore.AddSync(os.Stdout)

	var cores []zapcore.Core

	// 生产环境：只记录错误和更高级别的日志到文件
	if env == "production" {
		errorFileSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logDir, "error.log"),
			MaxSize:    100, // MB
			MaxBackups: 3,
			MaxAge:     30, // days
			Compress:   true,
		})
		// 只记录错误及以上级别到error.log
		cores = append(cores, zapcore.NewCore(fileEncoder, errorFileSyncer, zap.NewAtomicLevelAt(zap.ErrorLevel)))
	} else {
		// 非生产环境：记录所有级别的日志
		allFileSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logDir, "all.log"),
			MaxSize:    100, // MB
			MaxBackups: 3,
			MaxAge:     30, // days
			Compress:   true,
		})
		cores = append(cores, zapcore.NewCore(fileEncoder, allFileSyncer, zap.NewAtomicLevelAt(zap.InfoLevel)))

		// 也可以单独记录警告到warn.log
		warnFileSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logDir, "warn.log"),
			MaxSize:    100, // MB
			MaxBackups: 3,
			MaxAge:     30, // days
			Compress:   true,
		})
		cores = append(cores, zapcore.NewCore(fileEncoder, warnFileSyncer, zap.NewAtomicLevelAt(zap.WarnLevel)))
	}

	// 添加控制台输出（仅在非生产环境）
	if env != "production" {
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleSyncer, zap.NewAtomicLevelAt(zap.InfoLevel)))
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger.Sugar()
}

// Info logs an info message
func Info(message string, fields ...interface{}) {
	SugaredLogger.Infow(message, fields...)
}

// Error logs an error message
func Error(message string, fields ...interface{}) {
	SugaredLogger.Errorw(message, fields...)
}

// Warn logs a warning message
func Warn(message string, fields ...interface{}) {
	SugaredLogger.Warnw(message, fields...)
}

// Debug logs a debug message
func Debug(message string, fields ...interface{}) {
	SugaredLogger.Debugw(message, fields...)
}

// Panic logs a panic message
func Panic(message string, fields ...interface{}) {
	SugaredLogger.Panicw(message, fields...)
}

// WithRequest adds request information to the log
func WithRequest(c *gin.Context, message string, fields ...interface{}) {
	requestFields := []interface{}{
		"request_id", c.GetString("trace_id"),
		"method", c.Request.Method,
		"url", c.Request.URL.String(),
		"client_ip", c.ClientIP(),
		"user_agent", c.GetHeader("User-Agent"),
		"content_type", c.GetHeader("Content-Type"),
	}

	allFields := append(requestFields, fields...)
	Error(message, allFields...)
}

// WithRequestWarn adds request information to the log as warning
func WithRequestWarn(c *gin.Context, message string, fields ...interface{}) {
	requestFields := []interface{}{
		"request_id", c.GetString("trace_id"),
		"method", c.Request.Method,
		"url", c.Request.URL.String(),
		"client_ip", c.ClientIP(),
		"user_agent", c.GetHeader("User-Agent"),
		"content_type", c.GetHeader("Content-Type"),
	}

	allFields := append(requestFields, fields...)
	Warn(message, allFields...)
}
