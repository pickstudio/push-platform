package chi

import (
	"net/http"

	"github.com/pickstudio/push-platform/constants"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	"github.com/pickstudio/push-platform/pkg/er"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// RenderSuccessWithRenderer default render method for chi.
func RenderSuccessWithRenderer(w http.ResponseWriter, r *http.Request, v Renderer) {
	if v == nil {
		w.WriteHeader(oapiv1.StatusNoContent)
		return
	}
	if err := v.Render(w, r); err != nil {
		_ = RenderError(w, r, errors.Wrap(err, "response render"))
	}
}

// RenderSuccess default render method for chi.
func RenderSuccess(w http.ResponseWriter, r *http.Request, v any) {
	if v == nil {
		w.WriteHeader(oapiv1.StatusNoContent)
		return
	}
	b, err := jsoniter.Marshal(v)
	if err != nil {
		_ = RenderError(w, r, err)
	}
	w.WriteHeader(oapiv1.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		_ = RenderError(w, r, err)
	}
}

// ErrResponse is the HTTP response that reporting an error.
type ErrResponse struct {
	HTTPStatusCode int       `json:"-"`
	Status         ErrStatus `json:"status"`
}

// ErrStatus is the HTTP response body that reporting an error.
type ErrStatus struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// Render default render method for chi.
func (e ErrResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.HTTPStatusCode)
	w.Header().Set(oapiv1.HeaderContentType, oapiv1.MIMEApplicationJSONCharsetUTF8)

	b, err := jsoniter.Marshal(e)
	if err != nil {
		return err
	}
	_, _ = w.Write(b)

	return nil
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) error {
	httpStatusCode := constants.HTTPErrorToStatusCode(er.GetNamedErr(err))
	e := &ErrResponse{
		HTTPStatusCode: httpStatusCode,
		Status: ErrStatus{
			Code:    http.StatusText(httpStatusCode),
			Message: err.Error(),
		},
	}
	return e.Render(w, r)
}

// OAPIErrorHandler renders a response when OAPI error has occurred.
func OAPIErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if oapiv1.IsErrorTypeOfOAPI(err) {
		log.Warn().
			Err(err).
			Int("http_status_code", http.StatusBadRequest).
			Str("error_code", "invalid_argument").
			Send()
		e := &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			Status: ErrStatus{
				Code:    "invalid_argument",
				Message: err.Error(),
			},
		}
		_ = e.Render(w, r)
	}
}
