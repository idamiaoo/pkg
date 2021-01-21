package tracing

import (
	"fmt"
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func setupTracing(serviceName string, hostPort string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  hostPort,
		},
		RPCMetrics: true,
	}
	tracer, closer, err := cfg.New(serviceName, config.Logger(jaeger.NullLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

// SetupTracing 初始化 tracing
func SetupTracing(serviceName string, hostPort string) (opentracing.Tracer, io.Closer) {
	return setupTracing(serviceName, hostPort)
}
