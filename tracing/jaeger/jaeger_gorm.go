package jaeger

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var orm *traceOrm

type traceOrm struct {
	Ctx        context.Context
	ActiveSpan map[string]opentracing.Span
}

// NewTraceOrm 创建tranceorm对象
func NewTraceOrm() *traceOrm {
	if orm == nil {
		orm = &traceOrm{
			Ctx:        context.Background(),
			ActiveSpan: make(map[string]opentracing.Span),
		}
	}
	return orm
}

// 设置全局context
func (j *traceOrm) SetCtx(ctx context.Context) *traceOrm {
	j.Ctx = ctx
	return j
}

// 初始化全局orm中间件
func (j *traceOrm) InitCallBack() {
	gorm.DefaultCallback.Query().Before("gorm:query").Register("before_query", j.before)
	gorm.DefaultCallback.Query().After("gorm:query").Register("after_query", j.after)

	gorm.DefaultCallback.Create().Before("gorm:create").Register("before_create", j.before)
	gorm.DefaultCallback.Create().After("gorm:create").Register("after_create", j.after)

	gorm.DefaultCallback.Update().Before("gorm:update").Register("before_update", j.before)
	gorm.DefaultCallback.Update().After("gorm:update").Register("after_update", j.after)

	gorm.DefaultCallback.Delete().Before("gorm:delete").Register("before_delete", j.before)
	gorm.DefaultCallback.Delete().After("gorm:delete").Register("after_delete", j.after)
}

// 预处理
func (j *traceOrm) before(scope *gorm.Scope) {
	span, _ := opentracing.StartSpanFromContext(
		j.Ctx,
		fmt.Sprintf("DB.%s.%s", scope.Dialect().CurrentDatabase(), scope.TableName()),
		opentracing.Tag{Key: string(ext.Component), Value: "DB"},
		opentracing.Tag{Key: "DBName", Value: scope.Dialect().CurrentDatabase()},
	)
	j.ActiveSpan[scope.InstanceID()] = span
}

// 后处理
func (j *traceOrm) after(scope *gorm.Scope) {
	span := j.ActiveSpan[scope.InstanceID()]
	defer func() {
		span.Finish()
		delete(j.ActiveSpan, scope.InstanceID())
	}()
	ext.Error.Set(span, scope.HasError())
	if scope.HasError() {
		span.LogKV("error", scope.DB().Error.Error())
	}
	span.LogKV("SQL", scope.SQL)
	span.LogKV("SQLVAR", scope.SQLVars)
}
