package config

import (
	"strings"
	"time"
	"fmt"

	"github.com/spf13/viper"
)

var globalConfig *Config

// ============= 配置结构体 =============

// Config 应用配置
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Pagination PaginationConfig `mapstructure:"pagination"`
	RabbitMQ   RabbitMQConfig   `mapstructure:"rabbitmq"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port" json:"port"`
	Mode         string        `mapstructure:"mode" json:"mode"`
	Env          string        `mapstructure:"env" json:"env"` // development, production, test
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host         string        `mapstructure:"host" json:"host"`
	Port         int           `mapstructure:"port" json:"port"`
	Username     string        `mapstructure:"username" json:"username"`
	Password     string        `mapstructure:"password" json:"-"` // 不序列化密码
	Database     string        `mapstructure:"database" json:"database"`
	Charset      string        `mapstructure:"charset" json:"charset"`
	Socket       string        `mapstructure:"socket" json:"socket"`
	MaxIdleConns int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns int           `mapstructure:"max_open_conns" json:"max_open_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime" json:"max_lifetime"`
	MaxIdleTime  time.Duration `mapstructure:"max_idle_time" json:"max_idle_time"`
	Timeout      time.Duration `mapstructure:"timeout" json:"timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `mapstructure:"host" json:"host"`
	Port         int           `mapstructure:"port" json:"port"`
	Password     string        `mapstructure:"password" json:"-"`
	DB           int           `mapstructure:"db" json:"db"`
	PoolSize     int           `mapstructure:"pool_size" json:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns" json:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries" json:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout" json:"pool_timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `mapstructure:"secret" json:"-"`
	ExpireHours        int    `mapstructure:"expire_hours" json:"expire_hours"`
	RefreshExpireHours int    `mapstructure:"refresh_expire_hours" json:"refresh_expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level" json:"level"`
	Format   string `mapstructure:"format" json:"format"`
	Output   string `mapstructure:"output" json:"output"`
	FilePath string `mapstructure:"file_path" json:"file_path"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled" json:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second" json:"requests_per_second"`
	Burst             int  `mapstructure:"burst" json:"burst"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" json:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods" json:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers" json:"allowed_headers"`
}

// PaginationConfig 分页配置
type PaginationConfig struct {
	DefaultSize int `mapstructure:"default_size" json:"default_size"`
	MaxSize     int `mapstructure:"max_size" json:"max_size"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host          string `mapstructure:"host" json:"host"`
	Port          int    `mapstructure:"port" json:"port"`
	Username      string `mapstructure:"username" json:"username"`
	Password      string `mapstructure:"password" json:"-"`
	VHost         string `mapstructure:"vhost" json:"vhost"`
	Exchange      string `mapstructure:"exchange" json:"exchange"`
	PrefetchCount int    `mapstructure:"prefetch_count" json:"prefetch_count"`
	RetryLimit    int    `mapstructure:"retry_limit" json:"retry_limit"`
	RetryDelayMs  int    `mapstructure:"retry_delay_ms" json:"retry_delay_ms"`
}

// ============= 加载函数 =============

// Load 加载配置
func Load() (*Config, error) {
	// 设置配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")  // 优先从 configs 目录读取
	viper.AddConfigPath(".")         // 备选从当前目录读取

	// 环境变量支持
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// 配置文件不存在，使用默认值
	}

	// 设置默认值
	setDefaults()

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	// 保存全局配置
	globalConfig = &config

	return &config, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// ============= 默认值设置 =============

func setDefaults() {
	// Server
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// MySQL
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.username", "root")
	viper.SetDefault("mysql.password", "root")
	viper.SetDefault("mysql.database", "eshop_db")
	viper.SetDefault("mysql.charset", "utf8mb4")
	viper.SetDefault("mysql.socket", "")
	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)
	viper.SetDefault("mysql.max_lifetime", "1h")
	viper.SetDefault("mysql.max_idle_time", "5m")
	viper.SetDefault("mysql.timeout", "10s")
	viper.SetDefault("mysql.read_timeout", "30s")
	viper.SetDefault("mysql.write_timeout", "30s")

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 2)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")
	viper.SetDefault("redis.pool_timeout", "4s")

	// JWT
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.refresh_expire_hours", 168)

	// Log
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.file_path", "logs/app.log")

	// Rate Limit
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_second", 100)
	viper.SetDefault("rate_limit.burst", 200)

	// CORS
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Authorization"})

	// Pagination
	viper.SetDefault("pagination.default_size", 10)
	viper.SetDefault("pagination.max_size", 100)

	// RabbitMQ
	viper.SetDefault("rabbitmq.host", "localhost")
	viper.SetDefault("rabbitmq.port", 5672)
	viper.SetDefault("rabbitmq.username", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.vhost", "/")
	viper.SetDefault("rabbitmq.exchange", "eshop.events")
	viper.SetDefault("rabbitmq.prefetch_count", 10)
	viper.SetDefault("rabbitmq.retry_limit", 3)
	viper.SetDefault("rabbitmq.retry_delay_ms", 5000)
}

// ============= 配置验证 =============

func validateConfig(cfg *Config) error {
	// 验证服务器配置
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	// 验证环境
	validEnvs := map[string]bool{"development": true, "production": true, "test": true}
	if !validEnvs[cfg.Server.Env] {
		return fmt.Errorf("invalid environment: %s", cfg.Server.Env)
	}

	// 验证MySQL配置
	if cfg.MySQL.Username == "" {
		return fmt.Errorf("mysql username cannot be empty")
	}
	if cfg.MySQL.Database == "" {
		return fmt.Errorf("mysql database cannot be empty")
	}
	if cfg.MySQL.MaxOpenConns <= 0 {
		return fmt.Errorf("mysql max_open_conns must be > 0")
	}
	if cfg.MySQL.MaxIdleConns > cfg.MySQL.MaxOpenConns {
		return fmt.Errorf("mysql max_idle_conns (%d) cannot exceed max_open_conns (%d)", 
			cfg.MySQL.MaxIdleConns, cfg.MySQL.MaxOpenConns)
	}

	// 验证JWT
	if cfg.JWT.Secret == "your-secret-key-change-in-production" && cfg.Server.Env == "production" {
		return fmt.Errorf("JWT secret must be changed in production environment")
	}
	if cfg.JWT.ExpireHours <= 0 {
		return fmt.Errorf("jwt expire_hours must be > 0")
	}

	// 验证分页
	if cfg.Pagination.DefaultSize <= 0 {
		return fmt.Errorf("pagination default_size must be > 0")
	}
	if cfg.Pagination.MaxSize < cfg.Pagination.DefaultSize {
		return fmt.Errorf("pagination max_size must be >= default_size")
	}

	return nil
}