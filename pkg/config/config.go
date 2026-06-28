package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

var globalConfig *Config

// PaginationConfig 分页配置
type PaginationConfig struct {
	DefaultSize int
	MaxSize     int
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	VHost         string `mapstructure:"vhost"`
	Exchange      string `mapstructure:"exchange"`
	PrefetchCount int    `mapstructure:"prefetch_count"`
	RetryLimit    int    `mapstructure:"retry_limit"`
	RetryDelayMs  int    `mapstructure:"retry_delay_ms"`
}

// Config 应用配置
type Config struct {
	Server     ServerConfig
	MySQL      MySQLConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Log        LogConfig
	RateLimit  RateLimitConfig
	CORS       CORSConfig
	Pagination PaginationConfig
	RabbitMQ   RabbitMQConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	Charset      string
	Socket       string
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

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string
	ExpireHours        int
	RefreshExpireHours int
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string
	Format   string
	Output   string
	FilePath string
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerSecond int
	Burst             int
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")
	viper.AddConfigPath("./")

	// 环境变量覆盖
	viper.AutomaticEnv()

	// 然后在 Load 函数中
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// 设置默认值
	setDefaults()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 保存到全局变量
	globalConfig = &config

	return &config, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// MySQL配置
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.username", "root")
	viper.SetDefault("mysql.password", "root")
	viper.SetDefault("mysql.database", "eshop_db")
	viper.SetDefault("mysql.charset", "utf8mb4")
	viper.SetDefault("mysql.socket", "")
	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)

	// Redis配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT配置
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.refresh_expire_hours", 168)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.file_path", "logs/app.log")

	// 限流配置
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_second", 100)
	viper.SetDefault("rate_limit.burst", 200)

	// CORS配置
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Authorization"})

	// 分页配置
	viper.SetDefault("pagination.default_size", 10)
	viper.SetDefault("pagination.max_size", 100)

	// RabbitMQ配置
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
