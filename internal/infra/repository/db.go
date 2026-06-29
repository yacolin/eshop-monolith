package repository

import (
	repoModels "eshop-monolith/internal/infra/repository/models"

	"eshop-monolith/internal/product"

	invRepos "eshop-monolith/internal/inventory/domain/repositories"

	userRepos "eshop-monolith/internal/user/domain/repositories"

	orderRepos "eshop-monolith/internal/order/domain/repositories"

	cartRepos "eshop-monolith/internal/cart/domain/repositories"

	paymentRepos "eshop-monolith/internal/payment/domain/repositories"

	addressRepos "eshop-monolith/internal/address/domain/repositories"

	"eshop-monolith/pkg/config"
	"fmt"
	"os"
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
// resolveSocket 解析 socket 路径：配置明确设置时直接使用，否则自动探测常见路径
func resolveSocket(cfgSocket string) string {
	if cfgSocket != "" {
		return cfgSocket
	}
	for _, p := range []string{
		"/tmp/mysql.sock",                // macOS Homebrew
		"/var/run/mysqld/mysqld.sock",   // Linux/WSL2 Debian/Ubuntu
		"/run/mysqld/mysqld.sock",       // Linux/WSL2 systemd
		"/var/lib/mysql/mysql.sock",     // Linux RPM/CentOS
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

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&repoModels.CategoryPO{},
		&repoModels.AttributePO{},
		&repoModels.AttributeValuePO{},
		&repoModels.ProductPO{},
		&repoModels.InventoryPO{},

		&repoModels.ProductCategoryPO{},
		&repoModels.CategoryAttributePO{},

		&repoModels.OrderPO{},
		&repoModels.OrderItemPO{},

		&repoModels.UserPO{},
		&repoModels.UserInfoPO{},
		&repoModels.PermissionPO{},
		&repoModels.RolePO{},
		&repoModels.UserIdentityPO{},
		&repoModels.UserRolePO{},
		&repoModels.RolePermissionPO{},
		&repoModels.AuthTokenPO{},
		&repoModels.LoginHistoryPO{},

		&repoModels.PaymentPO{},
		&repoModels.PaymentMethodPO{},
		&repoModels.PaymentTransactionPO{},
		&repoModels.RefundPO{},

		&repoModels.CartPO{},
		&repoModels.CartItemPO{},

		&repoModels.CouponPO{},
		&repoModels.UserCouponPO{},
		&repoModels.PromotionPO{},
		&repoModels.PromotionProductPO{},

			&repoModels.SkuPO{},
			&repoModels.SkuAttributePO{},
			&repoModels.ProductAttributeValuePO{},

		&repoModels.NotificationPO{},

		&repoModels.AddressPO{},
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
	Redis        *redis.Client
	Brand        product.IbrandRepository
	Inventory    invRepos.IinventoryRepository
	Product      invRepos.IproductRepository
	Category     invRepos.IcategoryRepository
	Order        orderRepos.IorderRepository
	User         userRepos.IuserRepository
	UserInfo     userRepos.IuserInfoRepository
	UserIdentity userRepos.IuserIdentityRepository
	AuthToken    userRepos.IauthTokenRepository
	LoginHistory userRepos.IloginHistoryRepository
	Role         userRepos.IroleRepository
	Permission   userRepos.IpermissionRepository
	Cart         cartRepos.IcartRepository
	Payment      paymentRepos.IPaymentRepository
	Sku                  invRepos.IskuRepository
	ProductAttribute     invRepos.IproductAttributeRepository
	Attribute            invRepos.IattributeRepository
	CategoryAttribute    invRepos.IcategoryAttributeRepository
	Address              addressRepos.IaddressRepository
}

// NewRepositories 创建仓储集合
func NewRepositories(db *gorm.DB, redisClient *redis.Client) *Repositories {
	return &Repositories{
		Redis:        redisClient,
		Brand:        product.NewBrandRepository(db),
		Inventory:    invRepos.NewInventoryRepository(db),
		Product:      invRepos.NewProductRepository(db),
		Category:     invRepos.NewCategoryRepository(db),
		Order:        orderRepos.NewOrderRepository(db),
		User:         userRepos.NewUserRepository(db),
		UserInfo:     userRepos.NewUserInfoRepository(db),
		UserIdentity: userRepos.NewUserIdentityRepository(db),
		AuthToken:    userRepos.NewAuthTokenRepository(db),
		LoginHistory: userRepos.NewLoginHistoryRepository(db),
		Role:         userRepos.NewRoleRepository(db),
		Permission:   userRepos.NewPermissionRepository(db),
		Cart:         cartRepos.NewCachedCartRepository(cartRepos.NewCartRepository(db), redisClient, db),
		Payment:      paymentRepos.NewPaymentRepository(db),
		Sku:                  invRepos.NewSkuRepository(db),
		ProductAttribute:     invRepos.NewProductAttributeRepository(db),
		Attribute:            invRepos.NewAttributeRepository(db),
		CategoryAttribute:    invRepos.NewCategoryAttributeRepository(db),
		Address:              addressRepos.NewAddressRepository(db),
	}
}
