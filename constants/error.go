package constants

import (
	"github.com/pkg/errors"

	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	"github.com/pickstudio/push-platform/pkg/er"
)

var (
	// client side error.
	ErrInvalidRequest = errors.New("error, invalid request")
	ErrRequestTimeout = errors.New("error, request timeout")

	// server side error.
	ErrConfig       = errors.New("configuration struct is null")
	ErrWordNotExist = errors.New("doesn't exist token")

	ErrInternalServer = errors.New("error, internal server error")

	ErrHeaderAuditInvalid = errors.New("header audit value is invalid format")
)

var httpErrToStatus = map[error]int{
	ErrInvalidRequest: oapiv1.StatusBadRequest,
	ErrRequestTimeout: oapiv1.StatusRequestTimeout,

	ErrWordNotExist: oapiv1.StatusNotFound,

	ErrInternalServer: oapiv1.StatusInternalServerError,
}

func HTTPErrorToStatusCode(err error) int {
	if err == nil {
		return oapiv1.StatusInternalServerError
	}
	if v, ok := httpErrToStatus[err]; ok {
		return v
	}

	sourceErr := er.GetSourceErr(err)
	if oapiv1.IsErrorTypeOfOAPI(sourceErr) {
		return oapiv1.StatusBadRequest
	}

	return oapiv1.StatusInternalServerError
}
