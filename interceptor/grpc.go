package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	mdKeyLessonSN = "md_lesson_sn"
)

// LessonSNClientInterceptor grpc 调用中传输LessonSN
func LessonSNClientInterceptor(contextKey interface{}) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		value := ctx.Value(contextKey)
		if sValue, ok := value.(string); ok {
			md.Set(mdKeyLessonSN, sValue)
		}
		return invoker(ctx, method, req, resp, cc, opts...)
	}
}
