package v1

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type PostPushResponse struct {
	Status        Status           `json:"_status"`
	FailedMessage []*FailedMessage `json:"failed_message"`
}

func (res *PostPushResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(StatusOK)
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)

	b, err := jsoniter.Marshal(res)
	if err != nil {
		return err
	}
	_, _ = w.Write(b)
	return nil
}
