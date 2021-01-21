package logger

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
)

type contextKey struct{}

var lessonSNKey = contextKey{}

// LessonSNFromContext 从context获取lessonSN
func LessonSNFromContext(e *log.Entry) log.Fields {
	if e.Context != nil {
		val := e.Context.Value(lessonSNKey)
		if lessonSN, ok := val.(string); ok {
			return log.Fields{
				"lesson_sn": lessonSN,
			}
		}
	}
	return nil
}

func TestWithContextFields(t *testing.T) {
	Init(WithFields(log.Fields{
		"svrId":  "carl",
		"sender": "smurfs",
	}), WithFieldsFunc(LessonSNFromContext))
	ctx := context.WithValue(context.Background(), lessonSNKey, "1234567890")
	log.WithContext(ctx).Info("hello")
}
