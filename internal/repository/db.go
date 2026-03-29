package repository

import (
	"eshop-monolith/internal/domain/category"
	"eshop-monolith/internal/domain/product"
	"eshop-monolith/internal/domain/shared"
	"eshop-monolith/internal/pkg/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/redis/go-redis/v9"
)

// Config 数据库配置
type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// InitDB 初始化数据库连接
func InitDB(cfg config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&category.Category{},
		&product.Product{},
		&shared.ProductCategory{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// InitRedis 初始化Redis连接
func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	return client, nil
}

// Repositories 仓储集合
type Repositories struct {
	Order     OrderRepository
	Inventory InventoryRepository
	User      UserRepository
	Product   ProductRepository
	Category  CategoryRepository
}

// NewRepositories 创建仓储集合
func NewRepositories(db *gorm.DB, redisClient *redis.Client) *Repositories {
	return &Repositories{
		Order:     NewOrderRepository(db),
		Inventory: NewInventoryRepository(db),
		User:      NewUserRepository(db),
		Product:   NewProductRepository(db),
		Category:  NewCategoryRepository(db),
	}
}
