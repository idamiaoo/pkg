package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HTTPHeaderAcceptEncoding  string = "Accept-Encoding"
	HTTPHeaderAuthorization          = "Authorization"
	HTTPHeaderContentEncoding        = "Content-Encoding"
	HTTPHeaderContentLength          = "Content-Length"
	HTTPHeaderContentMD5             = "Content-MD5"
	HTTPHeaderContentType            = "Content-Type"
	HTTPHeaderContentLanguage        = "Content-Language"
	HTTPHeaderDate                   = "Date"
	HTTPHeaderExpires                = "Expires"
	HTTPHeaderHost                   = "Host"
	HTTPHeaderRange                  = "Range"
	HTTPHeaderLocation               = "Location"
	HTTPHeaderOrigin                 = "Origin"
	HTTPHeaderServer                 = "Server"
	HTTPHeaderUserAgent              = "User-Agent"
	HTTPHeaderNonce                  = "Nonce"
)

// Target .
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

func parseTarget(endpoint string) (*Target, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme == "" {
		if u, err = url.Parse("https://" + endpoint); err != nil {
			return nil, err
		}
	}
	target := &Target{Scheme: u.Scheme, Authority: u.Host}
	fmt.Println(u.Path)
	if len(u.Path) > 1 {
		target.Endpoint = u.Path[1:]
	}
	return target, nil
}

// Client .
type Client struct {
	opt    *options
	target *Target
	cc     *http.Client
}

// ClientOption .
type ClientOption func(*options)

// NewClient .
func NewClient(endpoint string, options ...ClientOption) (*Client, error) {
	opt := getDefaultOptions()

	for _, option := range options {
		option(opt)
	}

	target, err := parseTarget(endpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{
		opt:    opt,
		target: target,
		cc: &http.Client{
			Timeout:   opt.timeout,
			Transport: opt.transport,
		},
	}
	return client, err
}

func (client *Client) Invoke(ctx context.Context, method, path string, args interface{}, reply interface{}, opts ...CallOption) error {
	request := &Request{
		ctx:         ctx,
		client:      client,
		url:         fmt.Sprintf("%s://%s%s", client.target.Scheme, client.target.Authority, path),
		Method:      strings.ToUpper(method),
		QueryParam:  make(url.Values),
		FormData:    make(url.Values),
		Header:      make(http.Header),
		ContentType: client.opt.contentType,
		Time:        time.Now(),
	}

	if args != nil {
		if method == http.MethodGet {
			v, err := client.opt.form.Marshal(args)
			if err != nil {
				return err
			}
			request.SetQueryParamsFromValues(v)
		} else {
			switch request.ContentType {
			case FORM:
				data, err := client.opt.form.Marshal(args)
				if err != nil {
					return err
				}
				request.SetFormDataFromValues(data)
			case JSON:
				data, err := client.opt.json.Marshal(args)
				if err != nil {
					return err
				}
				request.Body = data
			}
		}
	}

	opts = combine(opts, client.opt.callOptions)
	// 执行中间件
	for _, opt := range opts {
		if err := opt.Before(request); err != nil {
			return err
		}
	}
	// 执行请求
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	// 执行中间件
	for i := len(opts) - 1; i >= 0; i-- {
		if err := opts[i].After(response); err != nil {
			return err
		}
	}

	if response.IsError() {
		return fmt.Errorf("status_code=%d", response.StatusCode())
	}

	if reply != nil {
		if err := client.opt.json.Unmarshal(response.Body(), reply); err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) Do(request *Request) (*Response, error) {
	if err := request.parseBody(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(request.Context(), request.Method, request.URL(), request.bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header = request.Header
	resp, err := client.cc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	response := &Response{
		request:     request,
		rawResponse: resp,
		time:        time.Now(),
	}
	response.body, _ = io.ReadAll(resp.Body)
	return response, nil
}
