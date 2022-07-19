package config

import (
	"time"

	"github.com/pickstudio/push-platform/const"
)

type Config struct {
	Env     string `env:"ENV,default=local"`   // for logging
	Debug   bool   `env:"DEBUG,default=false"` // for logging
	Release string // injected when it Dockerized by git hash

	LocalhostHttp struct {
		DSN     string        `env:"LOCALHOST_HTTP_DSN,default=0.0.0.0:50100" json:"localhost_http_dsn"`
		Timeout time.Duration `env:"LOCALHOST_HTTP_TIMEOUT,default=2s" json:"localhost_http_timeout"`
	}
}

func (c *Config) Self() (*Config, error) {
	if c == nil {
		return nil, _const.ErrConfig
	}
	return c, nil
}
