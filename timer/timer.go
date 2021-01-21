package timer

import (
	"context"
	"time"
)

// Timer 定时任务
type Timer interface {
	AddTask(ctx context.Context, delay time.Duration, task *Task) error
	WatchTask(ctx context.Context, taskName string, process func([]byte) error) error
}

// Task 任务信息
type Task struct {
	Name  string `json:"name" form:"name"`
	ID    string `json:"id" form:"id"`
	Value []byte `json:"value" form:"value"`
}
