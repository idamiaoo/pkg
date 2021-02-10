package logger

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestWithTraceID(t *testing.T) {
	Init(WithHooks(NewSenderHook("test")))
	ctx := context.Background()
	log.WithContext(ctx).Info("hello")
}
