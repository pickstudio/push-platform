package message

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-xray-sdk-go/xray"
	_const "github.com/pickstudio/push-platform/const"
	"github.com/pickstudio/push-platform/pkg/er"
	"time"
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

func (c *client) SendMessageBatch(ctx context.Context, list []types.SendMessageBatchRequestEntry) (int, []types.BatchResultErrorEntry, error) {
	op := er.GetOperator()
	ctx, root := xray.BeginSegment(ctx, _const.KeyProject+"-sqs")
	defer root.Close(nil)

	out, err := c.q.SendMessageBatch(
		ctx,
		&sqs.SendMessageBatchInput{
			Entries:  list,
			QueueUrl: aws.String(c.qDsn),
		})
	if err != nil {
		return 0, nil, er.WrapOp(err, op)
	}

	return len(out.Successful), out.Failed, nil
}

func (c *client) GetQueueAttributes(ctx context.Context) (*sqs.GetQueueAttributesOutput, error) {
	op := er.GetOperator()
	ctx, root := xray.BeginSegment(ctx, _const.KeyProject+"-sqs")
	defer root.Close(nil)

	out, err := c.q.GetQueueAttributes(
		ctx,
		&sqs.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:       aws.String(c.qDsn),
		},
	)
	if err != nil {
		return nil, er.WrapOp(err, op)
	}
	return out, nil
}
