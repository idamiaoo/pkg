package etcd

import (
	"context"
	"runtime/debug"
	"strings"
	"time"

	"github.com/katakurin/pkg/timer"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
)

const (
	timerPrefix = "/timer/"
	lockPrefix  = "/locker/"
	valuePrefix = "/taskvalue/"
)

type etcdTimer struct {
	client *v3.Client
	config v3.Config
}

func New(c v3.Config) (timer.Timer, error) {
	cli, err := v3.New(c)
	if err != nil {
		return nil, err
	}
	return &etcdTimer{
		client: cli,
		config: c,
	}, nil
}

// AddTask 添加定时任务
func (e *etcdTimer) AddTask(ctx context.Context, delay time.Duration, task *timer.Task) error {
	lease := v3.NewLease(e.client)
	grant, err := lease.Grant(ctx, int64(delay.Seconds()))
	_, err = e.client.KV.Put(ctx, timerPrefix+task.Name+task.ID, string(task.Value), v3.WithLease(grant.ID))
	if err != nil {
		return err
	}
	_, err = e.client.KV.Put(ctx, valuePrefix+task.Name+task.ID, string(task.Value))
	return err
}

// WatchTask 订阅定时任务
func (e *etcdTimer) WatchTask(ctx context.Context, taskName string, process func(task []byte) error) error {
	client, err := v3.New(e.config)
	if err != nil {
		return err
	}
	go func(ctx context.Context, client *v3.Client, taskName string, process func(task []byte) error) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
			log.Warn("任务监听异常退出")
			e.loop(ctx, client, taskName, process)
		}()
		e.loop(ctx, client, taskName, process)
	}(ctx, client, taskName, process)
	return nil
}

func (e *etcdTimer) loop(ctx context.Context, client *v3.Client, taskName string, process func([]byte) error) {
	log.WithField("task_name", taskName).Info("开始监听任务")
	watcher := v3.NewWatcher(client)
	ch := watcher.Watch(ctx, timerPrefix+taskName, v3.WithPrefix(), v3.WithRev(0), v3.WithFilterPut())
	for {
		select {
		case value := <-ch:
			for _, event := range value.Events {
				if event.Type == mvccpb.DELETE {
					taskID := strings.TrimPrefix(string(event.Kv.Key), timerPrefix+taskName)
					lg := log.WithFields(log.Fields{
						"task_name": taskName,
						"task_id":   taskID,
					})
					lg.WithField("value", event.Kv.Value).Info("监听到新的定时任务")
					lease, err := client.Grant(ctx, 60)
					if err != nil {
						lg.WithError(err).Error("生成租约失败")
						continue
					}
					lockerKey := lockPrefix + taskName + taskID
					cmp := v3.Compare(v3.CreateRevision(lockerKey), "=", 0)
					put := v3.OpPut(lockerKey, "", v3.WithLease(lease.ID))
					resp, err := client.Txn(ctx).If(cmp).Then(put).Commit()
					if err != nil {
						lg.WithError(err).Error("抢占位置失败")
						continue
					}
					if !resp.Succeeded {
						lg.Info("任务已在别处执行")
						continue
					}
					valueKey := valuePrefix + taskName + taskID
					getResp, err := client.Get(ctx, valueKey)
					if err != nil {
						lg.WithError(err).Error("获取任务信息失败")
						continue
					}
					if len(getResp.Kvs) != 1 {
						lg.WithError(err).Error("获取任务信息失败")
						continue
					}
					if err := process(getResp.Kvs[0].Value); err != nil {
						lg.WithError(err).Error("任务处理出错")
					}
					lg.Info("定时任务执行完成")
					_, err = client.Delete(ctx, valueKey)
					if err != nil {
						lg.WithError(err).Error("删除任务信息失败")
					}
				}
			}
		}
	}
}
