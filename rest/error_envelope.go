package rest

import "encoding/json"

type ErrorEnvelope struct {
	Error *Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}
