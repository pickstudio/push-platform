package handler

import (
	"net/http"

	edgechi "github.com/pickstudio/push-platform/edge/chi"
)

type handler struct {
}

func New() *handler {
	return &handler{}
}

// PostEnqueueFromDeadQueue DLQ로 빠진 실패한 에러메세지들을 다시 queue에다가 집어 넣을 수 있도록 합니다.
// (POST /_enqueue_from_dead_queue)
func (h *handler) PostEnqueueFromDeadQueue(w http.ResponseWriter, r *http.Request) {
	edgechi.RenderSuccess(w, r, map[string]string{
		"result": "_enqueue_from_dead_queue",
	})
}

// PostPush send push message immediately
// (POST /_push)
func (h *handler) PostPush(w http.ResponseWriter, r *http.Request) {
	edgechi.RenderSuccess(w, r, map[string]string{
		"result": "_push",
	})
}

// PostStatus send push message immediately
// (POST /_status)
func (h *handler) PostStatus(w http.ResponseWriter, r *http.Request) {
	edgechi.RenderSuccess(w, r, map[string]string{
		"result": "_status",
	})

}
