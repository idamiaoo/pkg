package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	ContentType ContentType
	Method      string
	QueryParam  url.Values
	FormData    url.Values
	Body        []byte
	Header      http.Header
	Time        time.Time
	url         string
	bodyReader  io.Reader
	ctx         context.Context
	client      *Client
}

func (r *Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) SetHeader(header, value string) *Request {
	r.Header.Set(header, value)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

func (r *Request) SetCookie(cookie *http.Cookie) *Request {
	r.Header.Add("Cookie", cookie.String())
	return r
}

func (r *Request) SetQueryParam(param, value string) *Request {
	r.QueryParam.Set(param, value)
	return r
}

func (r *Request) SetQueryParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}

func (r *Request) SetQueryParamsFromValues(params url.Values) *Request {
	for p, v := range params {
		for _, pv := range v {
			r.QueryParam.Add(p, pv)
		}
	}
	return r
}

func (r *Request) SetFormData(data map[string]string) *Request {
	for k, v := range data {
		r.FormData.Set(k, v)
	}
	return r
}

func (r *Request) SetFormDataFromValues(data url.Values) *Request {
	for k, v := range data {
		for _, kv := range v {
			r.FormData.Add(k, kv)
		}
	}
	return r
}

func (r *Request) SetFormDataFromStruct(data interface{}) *Request {
	values, _ := r.client.opt.form.Marshal(data)
	r.SetFormDataFromValues(values)
	return r
}

func (r *Request) SetBody(body []byte) *Request {
	r.Body = body
	return r
}

func (r *Request) URL() string {
	reqURL, err := url.Parse(r.url)
	if err != nil {
		return r.url
	}
	if len(r.QueryParam) > 0 {
		if len(strings.TrimSpace(reqURL.RawQuery)) == 0 {
			reqURL.RawQuery = r.QueryParam.Encode()
		} else {
			reqURL.RawQuery = reqURL.RawQuery + "&" + r.QueryParam.Encode()
		}
	}
	return reqURL.String()
}

func (r *Request) parseBody() (err error) {
	if r.Method == http.MethodPost {
		r.Header.Set(HTTPHeaderContentType, r.ContentType)
		if len(r.FormData) > 0 {
			r.Body = []byte(r.FormData.Encode())
			r.Header.Set(HTTPHeaderContentType, FORM)
		}
		if len(r.Body) > 0 {
			r.bodyReader = bytes.NewReader(r.Body)
		}
	}
	return
}
