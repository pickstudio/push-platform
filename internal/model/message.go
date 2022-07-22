package model

import (
	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	"github.com/pkg/errors"
)

var (
	ErrParsingMessageIdIsNotValid        = errors.New("error parsing 'message', 'id' is not valid")
	ErrParsingMessageUserIdIsNotValid    = errors.New("error parsing 'message', 'user_id' is not valid")
	ErrParsingMessageFromIsNotValid      = errors.New("error parsing 'message', 'from' is not valid")
	ErrParsingMessageServiceIsNotValid   = errors.New("error parsing 'message', 'service' is not valid")
	ErrParsingMessageDeviceIsNotValid    = errors.New("error parsing 'message', 'device' is not valid")
	ErrParsingMessagePushTokenIsNotValid = errors.New("error parsing 'message', 'push_token' is not valid")
	ErrParsingMessageViewTypeIsNotValid  = errors.New("error parsing 'message', 'view_type' is not valid")
	ErrParsingMessageViewIsNotValid      = errors.New("error parsing 'message', 'view' is not valid")
)

type Message struct {
	//Id  message identifier
	Id string `json:"id,omitempty"`
	// From for tracking who is send
	From string `json:"from,omitempty"`

	//Service one of service from pickstudio
	Service MessageService `json:"service,omitempty"`
	// UserId message owner
	UserId string `json:"user_id,omitempty"`

	// Device type to receive push message
	Device MessageDevice `json:"device,omitempty"`
	//PushToken actual push token by service
	PushToken string `json:"push_token,omitempty"`

	//ViewType view type of push message
	ViewType MessageViewType `json:"view_type,omitempty"`
	//PlainView view object is actual push message format
	PlainView *PlainView `json:"view,omitempty"`
}

func (m *Message) ToOAPI() *oapiv1.Message {
	var view any
	if m.ViewType == MessageViewTypePlain {
		view = m.PlainView
	}
	return &oapiv1.Message{
		Id:   m.Id,
		From: m.From,

		Service: m.Service.ToOAPI(),
		UserId:  m.UserId,

		Device:    m.Device.ToOAPI(),
		PushToken: m.PushToken,

		ViewType: m.ViewType.ToOAPI(),
		View:     view,
	}
}

func ParseOAPIMessage(v *oapiv1.Message) (*Message, error) {
	if v.Id == "" {
		return nil, ErrParsingMessageIdIsNotValid
	}
	if v.From == "" {
		return nil, ErrParsingMessageFromIsNotValid
	}

	service := ParseOAPIMessageService(v.Service)
	if service == MessageServiceUnknown {
		return nil, ErrParsingMessageServiceIsNotValid
	}
	if v.UserId == "" {
		return nil, ErrParsingMessageUserIdIsNotValid
	}

	device := ParseOAPIMessageDevice(v.Device)
	if device == MessageDeviceUnknown {
		return nil, ErrParsingMessageDeviceIsNotValid
	}
	if v.PushToken == "" {
		return nil, ErrParsingMessagePushTokenIsNotValid
	}

	viewType := ParseOAPIMessageViewType(v.ViewType)
	if viewType == MessageViewTypeUnknown {
		return nil, ErrParsingMessageViewTypeIsNotValid
	}
	if v.View == nil {
		return nil, ErrParsingMessageViewIsNotValid
	}

	m := &Message{
		Id:   v.Id,
		From: v.From,

		Service: service,
		UserId:  v.UserId,

		Device:    device,
		PushToken: v.PushToken,

		ViewType: viewType,
	}

	if viewType == MessageViewTypePlain {
		if view, ok := (v.View).(oapiv1.PlainView); ok {
			var err error
			m.PlainView, err = ParseOAPIPlainView(&view)
			if err != nil {
				return nil, err
			}
		}
	}
	return m, nil
}
