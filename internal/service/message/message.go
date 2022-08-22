package message

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/pickstudio/push-platform/constants"
	"github.com/pickstudio/push-platform/internal/model"
	"github.com/pickstudio/push-platform/pkg/arrays"
	"github.com/pickstudio/push-platform/pkg/er"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	ErrParsingJSON                  = errors.New("service.message: failed parsing json marshal/unmarshal")
	ErrInvalidParameterValueFromAWS = errors.New("service.message: Invalid parameter value about aws sqs service")
)

type messageAdapter interface {
	SendMessageBatch(ctx context.Context, list []types.SendMessageBatchRequestEntry, toDLQ bool) (int, []types.BatchResultErrorEntry, error)
	GetQueueAttributes(ctx context.Context, toDLQ bool) (*sqs.GetQueueAttributesOutput, error)
	ReceiveMessage(ctx context.Context, toDLQ bool) (*sqs.ReceiveMessageOutput, error)
}

type service struct {
	messageAdapter messageAdapter
}

func New(messageAdapter messageAdapter) *service {
	return &service{
		messageAdapter: messageAdapter,
	}
}

func (s *service) PushMessage(ctx context.Context, list []*model.Message) (int, []*model.FailedMessage, error) {
	op := er.GetOperator()

	var msgMap = map[string]*model.Message{}
	arrays.ForEach(list, func(msg *model.Message, _ int) {
		msgMap[msg.ID] = msg
	})

	var failedMessages []*model.FailedMessage
	resIn := arrays.Map(list, func(v *model.Message, _ int) *types.SendMessageBatchRequestEntry {
		var err error
		var body string
		if body, err = jsoniter.MarshalToString(v); err != nil {
			failedMessages = append(failedMessages, model.MessageWithError(v, err))
			return nil
		}

		return &types.SendMessageBatchRequestEntry{
			Id:          aws.String(v.ID),
			MessageBody: aws.String(body),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"From": {
					DataType:    aws.String("String"),
					StringValue: aws.String(constants.ValueProject),
				},
				"Service": {
					DataType:    aws.String("String"),
					StringValue: aws.String(string(v.Service)),
				},
				"Device": {
					DataType:    aws.String("String"),
					StringValue: aws.String(string(v.Device)),
				},
			},
			// MessageSystemAttributes: map[string]types.MessageSystemAttributeValue{
			//	"AWSTraceHeader": {
			//		DataType: aws.String("String"),
			//		//StringValue: aws.String(xray.GetSegment(ctx).TraceID),
			//	},
			// }, // .
		}
	})

	resIn = arrays.Filter(resIn, func(v *types.SendMessageBatchRequestEntry, _ int) bool {
		return v != nil
	})

	successCount, failedEntry, err := s.messageAdapter.SendMessageBatch(
		ctx,
		arrays.Map(resIn, func(v *types.SendMessageBatchRequestEntry, _ int) types.SendMessageBatchRequestEntry {
			return *v
		}),
		false,
	)
	if err != nil {
		err = er.WithNamedErr(er.WrapOp(err, op), constants.ErrInternalServer)
		return 0,
			arrays.Map(list, func(v *model.Message, i int) *model.FailedMessage {
				return model.MessageWithError(v, err)
			}),
			err
	}

	failedMessage := arrays.Map(failedEntry, func(v types.BatchResultErrorEntry, _ int) *model.FailedMessage {
		err = errors.WithMessagef(ErrInvalidParameterValueFromAWS, `id => %s
SenderFault => %v
Code => %s
Message => %s
`, *v.Id, v.SenderFault, *v.Code, *v.Message)

		m, ok := msgMap[aws.ToString(v.Id)]
		if !ok {
			return &model.FailedMessage{
				Error: err.Error(),
			}
		}

		return &model.FailedMessage{
			Message: m,
			Error:   err.Error(),
		}
	})
	failedMessage = arrays.Filter(failedMessage, func(v *model.FailedMessage, _ int) bool {
		return v != nil
	})
	if len(failedMessages) > 0 {
		s.pushMessageToDLQ(ctx, failedMessages)
	}
	return successCount, failedMessage, nil
}

func (s *service) ReceiveMessage(ctx context.Context) ([]*model.Message, error) {
	op := er.GetOperator()

	out, err := s.messageAdapter.ReceiveMessage(ctx, false)
	if err != nil {
		err = errors.WithMessage(err, op)
		return nil, err
	}

	var failedMessages []*model.FailedMessage
	res := arrays.Map(out.Messages, func(msg types.Message, _ int) *model.Message {
		var msgBody model.Message
		err = jsoniter.UnmarshalFromString(aws.ToString(msg.Body), &msgBody)
		if err != nil {
			err = errors.Errorf("%v/%s", err, aws.ToString(msg.Body))
			failedMessages = append(failedMessages, model.MessageWithError(nil, err))
		}
		return &msgBody
	})
	if len(failedMessages) > 0 {
		s.pushMessageToDLQ(ctx, failedMessages)
	}

	return res, nil
}

func (s *service) Status(ctx context.Context) (*model.Status, error) {
	op := er.GetOperator()

	_, err := s.messageAdapter.GetQueueAttributes(ctx, false)
	if err != nil {
		return nil, er.WrapOp(err, op)
	}
	return &model.Status{}, nil
}

func (s *service) pushMessageToDLQ(ctx context.Context, list []*model.FailedMessage) {
	op := er.GetOperator()
	var failedMessages []*model.FailedMessage

	resIn := arrays.Map(list, func(v *model.FailedMessage, _ int) *types.SendMessageBatchRequestEntry {
		var err error
		var body string

		if body, err = jsoniter.MarshalToString(v); err != nil {
			log.Error().Err(err).Interface("failed_message", v).Send()
			failedMessages = append(failedMessages, v)
			return nil
		}

		if v.Message == nil {
			failedMessages = append(failedMessages, v)
			return nil
		}

		return &types.SendMessageBatchRequestEntry{
			Id:                     aws.String(v.Message.ID),
			MessageDeduplicationId: aws.String(v.Message.ID),
			MessageGroupId:         aws.String(string(v.Message.Service)),
			MessageBody:            aws.String(body),
		}
	})
	resIn = arrays.Filter(resIn, func(v *types.SendMessageBatchRequestEntry, _ int) bool {
		return v != nil
	})
	log.Error().Interface("failed_messages", failedMessages).Str("operator", op).Send()

	_, _, err := s.messageAdapter.SendMessageBatch(
		ctx,
		arrays.Map(resIn, func(v *types.SendMessageBatchRequestEntry, _ int) types.SendMessageBatchRequestEntry {
			return *v
		}),
		true,
	)
	if err != nil {
		log.Error().Str("operator", op).Err(err).Send()
	}
}
