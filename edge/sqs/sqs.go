package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

// New just create SQS Client, with monitoring.
func New(ctx context.Context) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	// Instrumenting AWS SDK v2.
	awsv2.AWSV2Instrumentor(&cfg.APIOptions)
	// Using the Config value, create the DynamoDB.

	q := sqs.NewFromConfig(cfg)
	return q, nil
}
