package repository

import (
	"context"
	"database/sql"
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/product"
	"eshop-monolith/internal/trade"
	"eshop-monolith/internal/user"
	"eshop-monolith/pkg/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/redis/go-redis/v9"
)

// ============= 配置结构体 =============

// MySQLConfig MySQL配置（扩展）
type MySQLConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// RedisConfig Redis配置（扩展）
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
}

// ============= 默认值 =============

const (
	defaultCharset      = "utf8mb4"
	defaultMaxIdleConns = 10
	defaultMaxOpenConns = 100
	defaultMaxLifetime  = time.Hour
	defaultMaxIdleTime  = 5 * time.Minute
	defaultTimeout      = 10 * time.Second
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
)

// ============= 数据库初始化 =============

// InitDB 初始化MySQL连接
func InitDB(cfg config.MySQLConfig) (*gorm.DB, error) {
	// 1. 设置默认值
	cfg = applyMySQLDefaults(cfg)

	// 2. 构建DSN
	dsn := buildMySQLDSN(cfg)

	// 3. 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		// 性能优化配置
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 4. 获取底层SQL DB（注意这里是 *sql.DB）
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 5. 配置连接池（传递 *sql.DB）
	configureConnectionPool(sqlDB, cfg)

	// 6. 验证连接（传递 *sql.DB）
	if err := pingDB(sqlDB); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 7. 启动健康检查（传递 *sql.DB）
	startHealthCheck(sqlDB)

	return db, nil
}

// applyMySQLDefaults 应用MySQL配置默认值
func applyMySQLDefaults(cfg config.MySQLConfig) config.MySQLConfig {
	if cfg.Charset == "" {
		cfg.Charset = defaultCharset
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = defaultMaxIdleConns
	}
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = defaultMaxOpenConns
	}
	if cfg.MaxLifetime <= 0 {
		cfg.MaxLifetime = defaultMaxLifetime
	}
	if cfg.MaxIdleTime <= 0 {
		cfg.MaxIdleTime = defaultMaxIdleTime
	}
	return cfg
}

// buildMySQLDSN 构建MySQL DSN
func buildMySQLDSN(cfg config.MySQLConfig) string {
	// 基础参数
	baseParams := fmt.Sprintf(
		"charset=%s&parseTime=True&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s",
		cfg.Charset,
		getDurationOrDefault(cfg.Timeout, defaultTimeout),
		getDurationOrDefault(cfg.ReadTimeout, defaultReadTimeout),
		getDurationOrDefault(cfg.WriteTimeout, defaultWriteTimeout),
	)

	// 使用TCP连接
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, baseParams)
}

// getDurationOrDefault 获取duration或返回默认值
func getDurationOrDefault(d, defaultVal time.Duration) time.Duration {
	if d <= 0 {
		return defaultVal
	}
	return d
}

// configureConnectionPool 配置连接池（接收 *sql.DB）
func configureConnectionPool(sqlDB *sql.DB, cfg config.MySQLConfig) {
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(getDurationOrDefault(cfg.MaxLifetime, defaultMaxLifetime))
	sqlDB.SetConnMaxIdleTime(getDurationOrDefault(cfg.MaxIdleTime, defaultMaxIdleTime))
}

// pingDB 验证数据库连接（接收 *sql.DB）
func pingDB(sqlDB *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}

// startHealthCheck 启动健康检查（接收 *sql.DB）
func startHealthCheck(sqlDB *sql.DB) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := sqlDB.PingContext(ctx); err != nil {
				// 这里应该使用logger记录错误
				// logger.Errorf("Database health check failed: %v", err)
				_ = err // 忽略错误，避免编译警告
			}
			cancel()
		}
	}()
}

// ============= Redis初始化 =============

// InitRedis 初始化Redis连接
func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	// 应用默认值
	cfg = applyRedisDefaults(cfg)

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  getDurationOrDefault(cfg.DialTimeout, 5*time.Second),
		ReadTimeout:  getDurationOrDefault(cfg.ReadTimeout, 3*time.Second),
		WriteTimeout: getDurationOrDefault(cfg.WriteTimeout, 3*time.Second),
		PoolTimeout:  getDurationOrDefault(cfg.PoolTimeout, 4*time.Second),
	})

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

// applyRedisDefaults 应用Redis配置默认值
func applyRedisDefaults(cfg config.RedisConfig) config.RedisConfig {
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}
	if cfg.MinIdleConns <= 0 {
		cfg.MinIdleConns = 2
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	return cfg
}

// ============= 仓储集合 =============

// Repositories 仓储集合
type Repositories struct {
	Redis        *redis.Client
	Brand        product.IbrandRepository
	Inventory    inventory.IinventoryRepository
	Order        trade.IorderRepository
	User         user.IuserRepository
	UserInfo     user.IuserInfoRepository
	LoginHistory user.IloginHistoryRepository
}

// NewRepositories 创建仓储集合
func NewRepositories(db *gorm.DB, redisClient *redis.Client) *Repositories {
	return &Repositories{
		Redis:        redisClient,
		Brand:        product.NewBrandRepository(db),
		Inventory:    inventory.NewInventoryRepository(db),
		Order:        trade.NewOrderRepository(db),
		User:         user.NewUserRepository(db),
		UserInfo:     user.NewUserInfoRepository(db),
		LoginHistory: user.NewLoginHistoryRepository(db),
	}
}

// Close 关闭所有连接（优雅关闭时调用）
func (r *Repositories) Close() error {
	if r.Redis != nil {
		return r.Redis.Close()
	}
	return nil
}