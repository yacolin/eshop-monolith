package rabbitmq

import "eshop-monolith/pkg/config"

type Config struct {
	Host          string
	Port          int
	Username      string
	Password      string
	VHost         string
	Exchange      string
	PrefetchCount int
	RetryLimit    int
	RetryDelayMs  int
}

func NewConfig(cfg *config.RabbitMQConfig) Config {
	return Config{
		Host:          cfg.Host,
		Port:          cfg.Port,
		Username:      cfg.Username,
		Password:      cfg.Password,
		VHost:         cfg.VHost,
		Exchange:      cfg.Exchange,
		PrefetchCount: cfg.PrefetchCount,
		RetryLimit:    cfg.RetryLimit,
		RetryDelayMs:  cfg.RetryDelayMs,
	}
}
