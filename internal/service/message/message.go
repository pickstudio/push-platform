package message

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	_const "github.com/pickstudio/push-platform/const"
	"github.com/pickstudio/push-platform/internal/model"
	"github.com/pickstudio/push-platform/pkg/arrays"
	"github.com/pickstudio/push-platform/pkg/er"
)

type messageAdapter interface {
	SendMessageBatch(ctx context.Context, list []types.SendMessageBatchRequestEntry) (int, []types.BatchResultErrorEntry, error)
	GetQueueAttributes(ctx context.Context) (*sqs.GetQueueAttributesOutput, error)
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

	var msgByUUID = map[string]*model.Message{}

	successCount, failedEntry, err := s.messageAdapter.SendMessageBatch(
		ctx,
		arrays.Map(list, func(v *model.Message, _ int) types.SendMessageBatchRequestEntry {
			id := uuid.New().String()
			msgByUUID[id] = v
			body, _ := jsoniter.MarshalToString(v)
			return types.SendMessageBatchRequestEntry{
				Id:                     aws.String(id),
				MessageDeduplicationId: aws.String(v.Id),
				MessageBody:            aws.String(body),
				MessageAttributes:      map[string]types.MessageAttributeValue{},
				MessageSystemAttributes: map[string]types.MessageSystemAttributeValue{
					"AWSTraceHeader": {
						DataType:    aws.String("String"),
						StringValue: aws.String(xray.GetSegment(ctx).TraceID),
					},
				},
			}
		}),
	)
	if err != nil {
		return 0, nil, er.WithNamedErr(er.WrapOp(err, op), _const.ErrInternalServer)
	}
	failedMessage := arrays.Map(failedEntry, func(v types.BatchResultErrorEntry, _ int) *model.FailedMessage {
		var m = &model.Message{}
		_ = jsoniter.UnmarshalFromString(*v.Message, m)
		return &model.FailedMessage{
			Message: m,
			Error: fmt.Sprintf(`Id => %s
SenderFault => %v
Code => %s
Message => %s
`, *v.Id, v.SenderFault, *v.Code, *v.Message),
		}
	})

	return successCount, failedMessage, nil
}

func (s *service) Status(ctx context.Context) (*model.Status, error) {

	return nil, nil
}
