package sqs

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/pickstudio/push-platform/pkg/er"
)

// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/sqs/receivelpmessage/
type client struct {
	q *sqs.Client

	qTimeout time.Duration
	qDsn     string
	qName    string

	dlqTimeout time.Duration
	dlqDsn     string
	dlqName    string
}

func New(ctx context.Context, qName string, qTimeout time.Duration, dlqName string, dlqTimeout time.Duration) (*client, error) {
	op := er.GetOperator()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	q := sqs.NewFromConfig(cfg)

	qGetUrl, err := q.GetQueueUrl(
		ctx,
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(qName),
		},
	)

	if err != nil {
		return nil, er.WrapOp(err, op)
	}

	dlqGetUrl, err := q.GetQueueUrl(
		ctx,
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(dlqName),
		},
	)
	if err != nil {
		return nil, er.WrapOp(err, op)
	}

	return &client{
		q:        q,
		qName:    qName,
		qTimeout: qTimeout,
		qDsn:     *qGetUrl.QueueUrl,

		dlqName:    dlqName,
		dlqTimeout: dlqTimeout,
		dlqDsn:     *dlqGetUrl.QueueUrl,
	}, nil
}
