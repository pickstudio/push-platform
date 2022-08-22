package model

import oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"

type FailedMessage struct {
	Message *Message `json:"message"`
	Error   string   `json:"error"`
}

func (m *FailedMessage) ToOAPI() *oapiv1.FailedMessage {
	return &oapiv1.FailedMessage{
		Message: *m.Message.ToOAPI(),
		Error:   m.Error,
	}
}

func MessageWithError(msg *Message, err error) *FailedMessage {
	return &FailedMessage{
		Message: msg,
		Error:   err.Error(),
	}
}
