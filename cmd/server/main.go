// @title Eshop Monolith API
// @version 1.0
// @description 电商单体应用 API
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/infra/router"
	"eshop-monolith/pkg/config"
	"eshop-monolith/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	log.Println("Loading config...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Config loaded successfully: Server port %d", cfg.Server.Port)

	// 设置 Gin 运行模式（debug/release/test）
	gin.SetMode(cfg.Server.Mode)

	// 日志已通过 init() 函数自动初始化
	log.Println("Logger initialized")

	// 初始化数据库
	log.Println("Initializing database...")
	db, err := repository.InitDB(cfg.MySQL)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	log.Println("Database initialized successfully")

	// 初始化Redis
	log.Println("Initializing Redis...")
	redisClient, err := repository.InitRedis(cfg.Redis)
	if err != nil {
		log.Printf("Failed to initialize Redis: %v", err)
		os.Exit(1)
	}
	log.Println("Redis initialized successfully")

	// 初始化仓储
	repos := repository.NewRepositories(db, redisClient)

	// 初始化 RabbitMQ 客户端
	mqClient := rabbitmq.NewClient(rabbitmq.NewConfig(&cfg.RabbitMQ))
	defer mqClient.Close()

	// 初始化路由
	router := router.SetupRouter(cfg, repos, db, mqClient)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on port %d...", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
