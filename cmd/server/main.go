// @title Eshop Monolith API
// @version 1.0
// @description 电商单体应用 API
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
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
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/infra/router"
	"eshop-monolith/pkg/config"
	"eshop-monolith/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	log.Println("Loading config...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Config loaded successfully: Server port %d, Environment: %s", 
		cfg.Server.Port, cfg.Server.Env)

	// 2. 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	// 3. 初始化数据库（含连接验证）
	log.Println("Initializing database...")
	db, err := repository.InitDB(cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Printf("Database initialized successfully (MaxOpenConns=%d, MaxIdleConns=%d)", 
		cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)

	// 4. 初始化Redis
	log.Println("Initializing Redis...")
	redisClient, err := repository.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Printf("Redis initialized successfully (PoolSize=%d)", cfg.Redis.PoolSize)

	// 6. 初始化仓储
	repos := repository.NewRepositories(db, redisClient)
	defer func() {
		if err := repos.Close(); err != nil {
			log.Printf("Failed to close repositories: %v", err)
		}
	}()

	// 7. 初始化RabbitMQ客户端
	mqClient := rabbitmq.NewClient(rabbitmq.NewConfig(&cfg.RabbitMQ))
	defer func() {
		if err := mqClient.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ client: %v", err)
		}
	}()
	log.Println("RabbitMQ client initialized")

	// 8. 设置路由
	r := router.SetupRouter(cfg, repos, db, mqClient)
	log.Println("Router setup completed")

	// 9. 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second, // 空闲连接超时
		MaxHeaderBytes: 1 << 20,         // 1MB
	}

	// 10. 启动服务器（非阻塞）
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%d (Environment: %s)", 
			cfg.Server.Port, cfg.Server.Env)
		log.Printf("📝 API Documentation: http://localhost:%d/swagger/index.html", 
			cfg.Server.Port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 11. 优雅关闭
	waitForShutdown(server, db)
}

// waitForShutdown 等待中断信号并优雅关闭
func waitForShutdown(server *http.Server, db *gorm.DB) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// 监听关闭信号
	sig := <-quit
	logger.Info("Shutting down server...", "signal", sig.String())

	// 设置关闭超时（给现有请求5秒完成时间）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	// 关闭数据库连接
	closeDatabase(db)

	// 关闭Redis连接（已在Repositories.Close中处理）
	logger.Info("Server shutdown completed")
}

// closeDatabase 关闭数据库连接
func closeDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get database connection", "error", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Error("Failed to close database", "error", err)
	} else {
		logger.Info("Database connection closed successfully")
	}
}