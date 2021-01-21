package jaeger

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var redisMiddle *middleOption

type middleOption struct {
	Ctx context.Context
}

func NewRedisMiddle(ctx context.Context) *middleOption {
	if redisMiddle == nil {
		redisMiddle = &middleOption{Ctx: ctx}
	}
	return redisMiddle
}

//redis v6以下 包装process中间件。 v7有hook处理
//用法
// 1、创建客户端后调用：redisClient.WrapProcess(jaeger.RedisMiddle)
// 2、每次使用连接时调用传入context：jaeger.NewRedisMiddle(ctx)
func RedisMiddle(old func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
	return func(cmd redis.Cmder) error {
		span, _ := opentracing.StartSpanFromContext(
			redisMiddle.Ctx,
			fmt.Sprintf("REDIS.%s", cmd.Name()),
			opentracing.Tag{"METHOD", cmd.Name()},
		)
		defer span.Finish()
		ext.Component.Set(span, "REDIS")
		span.LogKV("cmd", cmd.Args())

		err := old(cmd)

		isError := cmd.Err() != nil && cmd.Err() != redis.Nil
		ext.Error.Set(span, isError)
		if isError {
			span.LogKV("error", cmd.Err().Error())
		}
		return err
	}
}
