package repository

import (
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/product"
	"eshop-monolith/internal/trade"
	"eshop-monolith/internal/user"

	"eshop-monolith/pkg/config"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/redis/go-redis/v9"
)

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

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

func resolveSocket(cfgSocket string) string {
	if cfgSocket != "" {
		return cfgSocket
	}
	for _, p := range []string{
		"/tmp/mysql.sock",
		"/var/run/mysqld/mysqld.sock",
		"/run/mysqld/mysqld.sock",
		"/var/lib/mysql/mysql.sock",
	} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func InitDB(cfg config.MySQLConfig) (*gorm.DB, error) {
	var dsn string
	sock := resolveSocket(cfg.Socket)
	if sock != "" {
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, sock, cfg.Database, cfg.Charset)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
	}

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
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	return client, nil
}

type Repositories struct {
	Redis        *redis.Client
	Brand        product.IbrandRepository
	Inventory    inventory.IinventoryRepository
	Order        trade.IorderRepository
	User         user.IuserRepository
	UserInfo     user.IuserInfoRepository
	LoginHistory user.IloginHistoryRepository
	Role         user.IroleRepository
	Permission   user.IpermissionRepository
}

func NewRepositories(db *gorm.DB, redisClient *redis.Client) *Repositories {
	return &Repositories{
		Redis:        redisClient,
		Brand:        product.NewBrandRepository(db),
		Inventory:    inventory.NewInventoryRepository(db),
		Order:        trade.NewOrderRepository(db),
		User:         user.NewUserRepository(db),
		UserInfo:     user.NewUserInfoRepository(db),
		LoginHistory: user.NewLoginHistoryRepository(db),
		Role:         user.NewRoleRepository(db),
		Permission:   user.NewPermissionRepository(db),
	}
}
