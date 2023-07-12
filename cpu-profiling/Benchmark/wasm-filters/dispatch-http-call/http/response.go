package http

import (
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

type Response struct {
	headers map[string]string
	body    string
}

func NewResponse(bodySize int) (*Response, error) {
	headersArray, err := proxywasm.GetHttpCallResponseHeaders()
	if err != nil {
		proxywasm.LogCritical("Failed to get http call response headers")
		return nil, err
	}
	headersMap := make(map[string]string)

	for _, item := range headersArray {
		headersMap[item[0]] = item[1]
	}

	body, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogCritical("Failed to get http call response body")
		return nil, err
	}

	return &Response{
		headers: headersMap,
		body:    string(body),
	}, nil
}

func (r *Response) GetStatus() int {
	code, _ := strconv.Atoi(r.headers[":status"])
	return code
}

func (r *Response) GetBody() string {
	return r.body
}

func (r *Response) GetHeaders() map[string]string {
	return r.headers
}

func (r *Response) GetHeader(key string) string {
	return r.headers[key]
}
