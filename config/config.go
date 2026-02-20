package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port        int    `env:"PORT" envDefault:"8080"`
	Mode        string `env:"APP_ENV" envDefault:"development"`
	AppName     string `env:"APP_NAME" envDefault:"search-service"`
	ServiceType string `env:"SERVICE_TYPE"` // http, grpc. event

	//Cache
	RedisAddr         string        `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword     string        `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB           int           `env:"REDIS_DB" envDefault:"0"`
	RequestLimiterTTL time.Duration `env:"REQUEST_LIMITER_TTL" envDefault:"10s"`
	RequestLimiterMax int64         `env:"REQUEST_LIMITER_MAX"`

	//API call
	MaxRetryCount int `env:"MAX_RETRY_COUNT" envDefault:"3"`
	RetryBackOff  int `env:"RETRY_BACKOFF" envDefault:"200"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load() (*Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{}
		if parseErr := env.Parse(cfg); parseErr != nil {
			err = parseErr
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}
