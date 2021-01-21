package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type config struct {
	fields          log.Fields // 默认字段
	fieldsFunc      []func(*log.Entry) log.Fields
	disableColors   bool
	fullTimestamp   bool
	reportCaller    bool
	level           log.Level
	timestampFormat string
}

type formatter struct {
	c  *config
	lf log.Formatter
}

// Format satisfies the logrus.Formatter interface.
func (f *formatter) Format(e *log.Entry) ([]byte, error) {
	for k, v := range f.c.fields {
		e.Data[k] = v
	}
	for _, f := range f.c.fieldsFunc {
		fields := f(e)
		for k, v := range fields {
			e.Data[k] = v
		}
	}
	return f.lf.Format(e)
}

// Option 日志选项
type Option func(*config)

// Init 初始化日志
func Init(opts ...Option) {
	c := &config{
		disableColors:   true,
		fullTimestamp:   true,
		reportCaller:    true,
		level:           log.DebugLevel,
		timestampFormat: "2006-01-02T15:04:05.000Z0700",
	}

	for _, opt := range opts {
		opt(c)
	}

	log.SetFormatter(&formatter{
		c: c,
		lf: &log.JSONFormatter{
			//DisableColors:   c.disableColors,
			//FullTimestamp:   c.fullTimestamp,
			TimestampFormat: c.timestampFormat,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				s := strings.Split(frame.Function, ".")
				funcname := s[len(s)-1]
				_, filename := path.Split(frame.File)
				filename = fmt.Sprintf("%s:%d", filename, frame.Line)
				return funcname, filename
			},
		},
	})
	log.SetReportCaller(c.reportCaller)
	log.SetLevel(c.level)
}

// WithFields 添加默认字段
func WithFields(fields log.Fields) Option {
	return func(c *config) {
		c.fields = fields
	}
}

// WithFieldsFunc 补充额外字段
func WithFieldsFunc(f func(*log.Entry) log.Fields) Option {
	return func(c *config) {
		c.fieldsFunc = append(c.fieldsFunc, f)
	}
}

// DisableColors 是否启动彩色模式
func DisableColors(v bool) Option {
	return func(c *config) {
		c.disableColors = v
	}
}

// Level 日志打印等级
func Level(level log.Level) Option {
	return func(c *config) {
		c.level = level
	}
}

// ReportCaller 是否添加调用信息字段
func ReportCaller(v bool) Option {
	return func(c *config) {
		c.reportCaller = v
	}
}

// TimestampFormat 设置日志时间格式
func TimestampFormat(v string) Option {
	return func(c *config) {
		c.timestampFormat = v
	}
}
