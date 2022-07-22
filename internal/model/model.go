package model

import (
	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
)

type MessageService string

const (
	MessageServiceUnknown    MessageService = "XX_UNKNOWN"
	MessageServiceBuddyStock MessageService = "BUDDYSTOCK"
	MessageServicePickMe     MessageService = "PICKME"
	MessageServiceDijkstra   MessageService = "DIJKSTRA"
)

func (v MessageService) ToOAPI() oapiv1.MessageService {
	switch v {
	case MessageServiceBuddyStock:
		return oapiv1.BUDDYSTOCK
	case MessageServicePickMe:
		return oapiv1.PICKME
	case MessageServiceDijkstra:
		return oapiv1.DIJKSTRA
	default:
		return oapiv1.DIJKSTRA
	}
}

func ParseOAPIMessageService(v oapiv1.MessageService) MessageService {
	switch v {
	case oapiv1.BUDDYSTOCK:
		return MessageServiceBuddyStock
	case oapiv1.PICKME:
		return MessageServicePickMe
	case oapiv1.DIJKSTRA:
		return MessageServiceDijkstra
	default:
		return MessageServiceUnknown
	}
}

// MessageDevice MessageDevice
type MessageDevice string

const (
	MessageDeviceUnknown       MessageDevice = "XX_UNKNOWN"
	MessageDeviceIOS           MessageDevice = "IOS"
	MessageDeviceAndroid       MessageDevice = "ANDROID"
	MessageDeviceDesktopChrome MessageDevice = "DESKTOP_CHROME"
)

func ParseOAPIMessageDevice(v oapiv1.MessageDevice) MessageDevice {
	switch v {
	case oapiv1.IOS:
		return MessageDeviceIOS
	case oapiv1.ANDROID:
		return MessageDeviceAndroid
	case oapiv1.DESKTOPCHROME:
		return MessageDeviceDesktopChrome
	default:
		return MessageDeviceUnknown
	}
}

func (v MessageDevice) ToOAPI() oapiv1.MessageDevice {
	switch v {
	case MessageDeviceIOS:
		return oapiv1.IOS
	case MessageDeviceAndroid:
		return oapiv1.ANDROID
	case MessageDeviceDesktopChrome:
		return oapiv1.DESKTOPCHROME
	default:
		return oapiv1.DESKTOPCHROME
	}
}

// MessageViewType MessageViewType
type MessageViewType string

const (
	MessageViewTypeUnknown MessageViewType = "XX_UNKNOWN"
	MessageViewTypePlain   MessageViewType = "PLAIN"
)

func ParseOAPIMessageViewType(v oapiv1.MessageViewType) MessageViewType {
	switch v {
	case oapiv1.PLAIN:
		return MessageViewTypePlain
	default:
		return MessageViewTypeUnknown
	}
}

func (v MessageViewType) ToOAPI() oapiv1.MessageViewType {
	switch v {
	case MessageViewTypePlain:
		return oapiv1.PLAIN
	default:
		return oapiv1.PLAIN
	}
}
