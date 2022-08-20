package message

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/google/uuid"
	"github.com/pickstudio/push-platform/constants"
	"github.com/pickstudio/push-platform/pkg/er"
)

type client struct {
	q *sqs.Client

	qTimeout time.Duration
	qDsn     string
	qName    string

	dlqTimeout time.Duration
	dlqDsn     string
	dlqName    string
}

func New(
	ctx context.Context, q *sqs.Client,
	qName string, qTimeout time.Duration,
	dlqName string, dlqTimeout time.Duration,
) (*client, error) {
	op := er.GetOperator()

	qURL, err := q.GetQueueUrl(
		ctx,
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(qName),
		},
	)

	if err != nil {
		return nil, er.WrapOp(err, op)
	}

	dlqURL, err := q.GetQueueUrl(
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
		qDsn:     *qURL.QueueUrl,

		dlqName:    dlqName,
		dlqTimeout: dlqTimeout,
		dlqDsn:     *dlqURL.QueueUrl,
	}, nil
}

func (c *client) SendMessageBatch(ctx context.Context, list []types.SendMessageBatchRequestEntry, toDlq bool) (int, []types.BatchResultErrorEntry, error) {
	op := er.GetOperator()
	ctx, root := xray.BeginSegment(ctx, constants.KeyProject+"-sqs")
	defer root.Close(nil)

	dsn := c.qDsn
	if toDlq {
		dsn = c.dlqDsn
	}

	out, err := c.q.SendMessageBatch(
		ctx,
		&sqs.SendMessageBatchInput{
			Entries:  list,
			QueueUrl: aws.String(dsn),
		})
	if err != nil {
		return 0, nil, er.WrapOp(err, op)
	}

	return len(out.Successful), out.Failed, nil
}

func (c *client) GetQueueAttributes(ctx context.Context, toDlq bool) (*sqs.GetQueueAttributesOutput, error) {
	op := er.GetOperator()
	ctx, root := xray.BeginSegment(ctx, constants.KeyProject+"-sqs")
	defer root.Close(nil)

	dsn := c.qDsn
	if toDlq {
		dsn = c.dlqDsn
	}

	out, err := c.q.GetQueueAttributes(
		ctx,
		&sqs.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:       aws.String(dsn),
		},
	)
	if err != nil {
		return nil, er.WrapOp(err, op)
	}

	return out, nil
}

func (c *client) ReceiveMessage(ctx context.Context, toDlq bool) (*sqs.ReceiveMessageOutput, error) {
	op := er.GetOperator()
	ctx, root := xray.BeginSegment(ctx, constants.KeyProject+"-sqs")
	defer root.Close(nil)

	dsn := c.qDsn
	if toDlq {
		dsn = c.dlqDsn
	}

	out, err := c.q.ReceiveMessage(
		ctx,
		&sqs.ReceiveMessageInput{
			AttributeNames:          []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:                aws.String(dsn),
			ReceiveRequestAttemptId: aws.String(uuid.NewString()),
			WaitTimeSeconds:         10,
			// MessageAttributeNames: nil,
			// MaxNumberOfMessages: 1, // 1 is default
			// VisibilityTimeout: 10, // .
		},
	)
	if err != nil {
		return nil, er.WrapOp(err, op)
	}

	return out, nil
}
