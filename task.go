package ziva

import "net/http"

type Task struct {
	Url string

	Method string

	Payload []byte

	FormData FormData

	Header *http.Header

	Data map[string]interface{}
}
