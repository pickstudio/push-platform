package config

import (
	"time"

	"github.com/pickstudio/push-platform/constants"
)

type Config struct {
	Env           string `env:"ENV,default=local"`   // for logging.
	Debug         bool   `env:"DEBUG,default=false"` // for logging.
	Release       string // injected when it Dockerized by git hash.
	Monitoring    bool   `env:"MONITORING,default=false"`
	LocalhostHTTP struct {
		DSN     string        `env:"LOCALHOST_HTTP_DSN,default=0.0.0.0:50100" json:"localhost_http_dsn"`
		Timeout time.Duration `env:"LOCALHOST_HTTP_TIMEOUT,default=2s" json:"localhost_http_timeout"`
	}

	AWSSQSQueue struct {
		Name    string        `env:"AWS_SQS_QUEUE_NAME" json:"aws_sqs_queue_name"`
		Timeout time.Duration `env:"AWS_SQS_QUEUE_TIMEOUT" json:"aws_sqs_queue_timeout"`
	}
	AWSSQSDeadLetterQueue struct {
		Name    string        `env:"AWS_SQS_DEADLETTER_QUEUE_NAME" json:"aws_sqs_deadletter_queue_name"`
		Timeout time.Duration `env:"AWS_SQS_DEADLETTER_QUEUE_TIMEOUT" json:"aws_sqs_deadletter_queue_timeout"`
	}
}

func (c *Config) Self() (*Config, error) {
	if c == nil {
		return nil, constants.ErrConfig
	}
	return c, nil
}
