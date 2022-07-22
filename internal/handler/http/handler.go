package http

import (
	"context"
	"encoding/json"
	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	"github.com/pickstudio/push-platform/internal/model"
	"github.com/pickstudio/push-platform/pkg/arrays"
	"github.com/pickstudio/push-platform/pkg/er"
	"net/http"

	edgechi "github.com/pickstudio/push-platform/edge/chi"
)

type messageService interface {
	PushMessage(ctx context.Context, list []*model.Message) (int, []*model.FailedMessage, error)
}

type handler struct {
	messageService messageService
}

func New(messageService messageService) *handler {
	return &handler{
		messageService: messageService,
	}
}

// PostPush send push message immediately
// (POST /_push)
func (h *handler) PostPush(w http.ResponseWriter, r *http.Request) {
	op := er.GetOperator()

	ctx := r.Context()
	var body oapiv1.PostPushJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		_ = edgechi.RenderError(w, r, err)
		return
	}
	var failedMessage []*oapiv1.FailedMessage
	list := arrays.Map(*body.Messages, func(v oapiv1.Message, i int) *model.Message {
		m, err := model.ParseOAPIMessage(&v)
		if err != nil {
			failedMessage = append(failedMessage, &oapiv1.FailedMessage{
				Error:   err.Error(),
				Message: v,
			})
			return nil
		}
		return m
	})
	list = arrays.Filter(list, func(v *model.Message, _ int) bool {
		return v != nil
	})

	_, fm, err := h.messageService.PushMessage(ctx, list)
	if err != nil {
		err := er.WrapOp(err, op)
		_ = edgechi.RenderError(w, r, err)
		return
	}
	if len(fm) > 0 {
		arrays.ForEach(fm, func(v *model.FailedMessage, _ int) {
			failedMessage = append(failedMessage, v.ToOAPI())
		})
	}

	resBody := &oapiv1.PostPushResponse{
		Status:        oapiv1.Status{},
		FailedMessage: failedMessage,
	}
	_ = resBody.Render(w, r)
	return
}

// PostStatus send push message immediately
// (POST /_status)
func (h *handler) PostStatus(w http.ResponseWriter, r *http.Request) {
	edgechi.RenderSuccess(w, r, map[string]string{
		"result": "_status",
	})
}

// PostEnqueueFromDeadQueue DLQ로 빠진 실패한 에러메세지들을 다시 queue에다가 집어 넣을 수 있도록 합니다.
// (POST /_enqueue_from_dead_queue)
func (h *handler) PostEnqueueFromDeadQueue(w http.ResponseWriter, r *http.Request) {
	edgechi.RenderSuccess(w, r, map[string]string{
		"result": "_enqueue_from_dead_queue",
	})
}
