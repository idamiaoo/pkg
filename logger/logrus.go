package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Caller 函数调用信息
func Caller(frame *runtime.Frame) (function string, file string) {
	s := strings.Split(frame.Function, ".")
	function = s[len(s)-1]
	_, file = path.Split(frame.File)
	file = fmt.Sprintf("%s:%d", file, frame.Line)
	return
}

type options struct {
	formatter    log.Formatter
	reportCaller bool
	level        log.Level
	hooks        []log.Hook
}

// Option 初始化选项
type Option interface {
	apply(*options)
}

type funcInitOption struct {
	f func(*options)
}

func (fio *funcInitOption) apply(io *options) {
	fio.f(io)
}

func newFuncInitOption(f func(*options)) *funcInitOption {
	return &funcInitOption{
		f: f,
	}
}

// Init 初始化日志
func Init(opts ...Option) {
	initOptions := &options{
		reportCaller: false,
		level:        log.DebugLevel,
	}

	for _, opt := range opts {
		opt.apply(initOptions)
	}
	if initOptions.formatter != nil {
		log.SetFormatter(initOptions.formatter)
	}

	if initOptions.reportCaller {
		log.SetReportCaller(initOptions.reportCaller)
	}

	log.SetLevel(initOptions.level)

	for _, hook := range initOptions.hooks {
		log.AddHook(hook)
	}
}

// WithLevel 日志打印等级
func WithLevel(level log.Level) Option {
	return newFuncInitOption(func(o *options) {
		o.level = level
	})
}

// ReportCaller 是否添加调用信息字段
func ReportCaller() Option {
	return newFuncInitOption(func(o *options) {
		o.reportCaller = true
	})
}

// WithHooks 添加 hooks
func WithHooks(hooks ...log.Hook) Option {
	return newFuncInitOption(func(o *options) {
		o.hooks = hooks
	})
}
