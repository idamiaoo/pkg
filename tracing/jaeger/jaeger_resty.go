package jaeger

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// RestyTrace resty请求追踪
type RestyTrace struct {
	span opentracing.Span
	logReq  bool
	logResp bool
}

// NewRestyTrace .
func NewRestyTrace() *RestyTrace {
	return &RestyTrace{}
}

type TraceOptFn func(*RestyTrace)

// LogReq 记录请求参数据
func LogReq() TraceOptFn {
	return func(opt *RestyTrace) {
		opt.logReq = true
	}
}
// LogResp 记录请求返回数据
func LogResp() TraceOptFn {
	return func(opt *RestyTrace) {
		opt.logResp = true
	}
}


// NewRestyClient .
func NewRestyClient(opts ...TraceOptFn) *resty.Client {
	trace := NewRestyTrace()
	for _, fn := range opts {
		fn(trace)
	}
	client := resty.New()
	client.OnBeforeRequest(trace.BeforeRequest())
	client.OnAfterResponse(trace.AfterResponse())
	return client
}

// BeforeRequest .
func (t *RestyTrace) BeforeRequest() resty.RequestMiddleware {
	return func(c *resty.Client, req *resty.Request) error {
		span, _ := opentracing.StartSpanFromContext(req.Context(), req.URL)
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
			return err
		}
		t.span = span
		ext.HTTPMethod.Set(span, req.Method)
		ext.HTTPUrl.Set(span, req.URL)
		ext.Component.Set(span, "HTTP")
		if t.logReq {
			if req.Method == resty.MethodGet {
				param, _ := json.Marshal(req.QueryParam)
				span.LogKV("http.param", string(param))
			}
			if req.Method == resty.MethodPost {
				param, _ := json.Marshal(req.Body)
				span.LogKV("http.param", string(param))
			}
		}
		return nil
	}
}

// AfterResponse .
func (t *RestyTrace) AfterResponse() resty.ResponseMiddleware {
	return func(c *resty.Client, resp *resty.Response) error {
		ext.HTTPStatusCode.Set(t.span, uint16(resp.StatusCode()))
		ext.Error.Set(t.span, resp.StatusCode() <= 199 || resp.StatusCode() >= 300)
		if t.logResp {
			t.span.LogKV("http.response", string(resp.Body()))
		}
		t.span.Finish()
		return nil
	}
}
