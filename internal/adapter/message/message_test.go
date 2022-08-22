package message

import (
	"context"
	"fmt"
	"testing"
	"time"

	edgesqs "github.com/pickstudio/push-platform/edge/sqs"

	"github.com/Netflix/go-env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/pickstudio/push-platform/internal/config"
)

var (
	cfg config.Config
)

func init() {
	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Interface("config", cfg).Msg("http_server start")
}

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := edgesqs.New(ctx)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}
	cc, err := New(
		ctx,
		c,
		cfg.AWSSQSQueue.Name,
		cfg.AWSSQSQueue.Timeout,
		cfg.AWSSQSDeadLetterQueue.Name,
		cfg.AWSSQSDeadLetterQueue.Timeout,
	)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}

	sendRes, err := c.SendMessage(
		ctx,
		&sqs.SendMessageInput{
			DelaySeconds: 10,
			MessageAttributes: map[string]types.MessageAttributeValue{
				"Title": {
					DataType:    aws.String("String"),
					StringValue: aws.String("The Whistler"),
				},
				"Author": {
					DataType:    aws.String("String"),
					StringValue: aws.String("John Grisham"),
				},
				"WeeksOn": {
					DataType:    aws.String("Number"),
					StringValue: aws.String("6"),
				},
			},
			MessageBody: aws.String("Information about the NY Times fiction bestseller for the week of 12/11/2016."),
			QueueUrl:    aws.String(cc.qDsn),
		},
	)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}
	fmt.Println(jsoniter.MarshalToString(sendRes))

	msgBody, err := cc.q.ReceiveMessage(
		ctx,
		&sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(cc.qDsn),
			AttributeNames:        []types.QueueAttributeName{"QueueArn"},
			MessageAttributeNames: []string{"Title"},
		},
	)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}
	fmt.Println(jsoniter.MarshalToString(msgBody))

	attr, err := c.GetQueueAttributes(
		ctx,
		&sqs.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:       aws.String(cc.qDsn),
		},
	)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}
	fmt.Println(jsoniter.MarshalToString(attr.Attributes))
	fmt.Println(jsoniter.MarshalToString(attr.ResultMetadata))
}
