package http

import "encoding/json"

type ResponseCode int

const (
	Success ResponseCode = iota
	InvalidRequest
)

type Response struct {
	Code  ResponseCode `json:"code"`
	Data  interface{}  `json:"data,omitempty"`
	Error string       `json:"error,omitempty"`
}

func (r *Response) Marshal() string {
	b, _ := json.Marshal(r)
	return string(b)
}
