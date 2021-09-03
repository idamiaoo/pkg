package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// options 客户端配置
type options struct {
	timeout     time.Duration
	transport   http.RoundTripper
	callOptions []CallOption
	form        FormCodec
	json        JsonCodec
	contentType ContentType
}

func getDefaultOptions() *options {
	opt := options{}
	opt.transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	opt.form = NewDefaultFormCodec()
	opt.json = NewDefaultJsonCodec()
	opt.contentType = JSON
	return &opt
}

// WithDefaultCallOptions 设置默认请求选项
func WithDefaultCallOptions(cos ...CallOption) ClientOption {
	return func(opt *options) {
		opt.callOptions = append(opt.callOptions, cos...)
	}
}

// WithContentType 设置post请求的content类型
func WithContentType(contentType string) ClientOption {
	return func(opt *options) {
		if contentType != "" {
			opt.contentType = contentType
		}
	}
}
