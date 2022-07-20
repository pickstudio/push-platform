package sqs

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

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
	cfg        config.Config
	httpServer *http.Server
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

	c, err := New(
		ctx,
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

	sendRes, err := c.q.SendMessage(
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
			QueueUrl:    aws.String(c.qDsn),
		},
	)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}
	fmt.Println(jsoniter.MarshalToString(sendRes))

	attr, err := c.q.GetQueueAttributes(
		ctx,
		&sqs.GetQueueAttributesInput{
			QueueUrl:       aws.String(c.qDsn),
			AttributeNames: []types.QueueAttributeName{"Title", "Author"},
		},
	)
	if err != nil {
		fmt.Println(err)
		assert.Empty(t, err)
		return
	}
	fmt.Println(jsoniter.MarshalToString(attr))

}
