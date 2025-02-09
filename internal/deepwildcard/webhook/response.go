package webhook

import (
	"encoding/json"
	"io"

	"github.com/smallstep/certificates/webhook"
)

type Response struct {
	webhook.ResponseBody
}

func CreateDenied(code string, message string) *Response {
	r := &Response{}
	r.Data = nil
	r.Allow = false
	r.Error = &webhook.Error{
		Message: message,
		Code:    code,
	}
	return r
}

func CreateAllowed() *Response {
	r := &Response{}
	r.Data = nil
	r.Allow = true
	r.Error = &webhook.Error{
		Message: "",
		Code:    "",
	}
	return r
}

func (r *Response) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}
