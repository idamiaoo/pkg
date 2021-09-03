package http

import (
	"net/http"
	"strings"
	"time"
)

type Response struct {
	request     *Request
	rawResponse *http.Response
	time        time.Time
	body        []byte
}

func (r *Response) Body() []byte {
	if r.body == nil {
		return []byte{}
	}
	return r.body
}

func (r *Response) Time() time.Time {
	return r.time
}

func (r *Response) Cost() time.Duration {
	return r.time.Sub(r.request.Time)
}

func (r *Response) StatusCode() int {
	if r.rawResponse == nil {
		return 0
	}
	return r.rawResponse.StatusCode
}

// String method returns the body of the server response as String.
func (r *Response) String() string {
	if r.body == nil {
		return ""
	}
	return strings.TrimSpace(string(r.body))
}

// IsSuccess method returns true if HTTP status `code >= 200 and <= 299` otherwise false.
func (r *Response) IsSuccess() bool {
	return r.StatusCode() > 199 && r.StatusCode() < 300
}

// IsError method returns true if HTTP status `code >= 400` otherwise false.
func (r *Response) IsError() bool {
	return r.StatusCode() > 399
}
