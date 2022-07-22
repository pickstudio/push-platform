package v1

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type PostPushResponse struct {
	Status        Status           `json:"_status"`
	FailedMessage []*FailedMessage `json:"failed_message"`
}

func (res *PostPushResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(StatusOK)
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)

	b, err := jsoniter.Marshal(res)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}