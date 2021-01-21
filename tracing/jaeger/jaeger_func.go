package jaeger

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// StartFuncSpan .
func StartFuncSpan(ctx context.Context) (opentracing.Span, context.Context) {
	funcName := funcName()
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		fmt.Sprintf("Func.%s", funcName),
		opentracing.Tag{Key: "Method", Value: funcName},
	)
	ext.Component.Set(span, "FUNC")
	return span, ctx
}

// funcName 获取正在运行的函数名
func funcName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	arr := strings.Split(f.Name(), "/")
	name := arr[len(arr)-1:][0]
	return name
}
