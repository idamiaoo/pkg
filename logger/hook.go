package logger

import (
	log "github.com/sirupsen/logrus"
)

type SenderHook struct {
	sender string
}

func (h *SenderHook) Levels() []log.Level {
	return log.AllLevels
}
func (h *SenderHook) Fire(e *log.Entry) error {
	e.Data["sender"] = h.sender
	return nil
}

// NewSenderHook 添加 sender 字段
func NewSenderHook(sender string) log.Hook {
	return &SenderHook{
		sender: sender,
	}
}

type ContextKey struct{}

// TraceIDKey  trace id key
var TraceIDKey = ContextKey{}

type ExtractFiledFormContextHook struct {
	key   ContextKey
	filed string
}

func (h *ExtractFiledFormContextHook) Levels() []log.Level {
	return log.AllLevels
}
func (h *ExtractFiledFormContextHook) Fire(e *log.Entry) error {
	if e.Context != nil {
		val := e.Context.Value(h.key)
		e.Data[h.filed] = val
	}
	return nil
}

// NewTraceIDFromContextHook 从context获取trace_id
func NewTraceIDFromContextHook(key ContextKey, filed string) log.Hook {
	return &ExtractFiledFormContextHook{
		key:   key,
		filed: filed,
	}
}
