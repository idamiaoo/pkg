package etcdlocker

import (
	"context"

	"github.com/katakurin/pkg/locker"
	"github.com/pkg/errors"
	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// 基于etcdv3的分布式锁
type etcd3Lock struct {
	id      string
	mutex   *concurrency.Mutex
	session *concurrency.Session
}

func Acquire(client *v3.Client, id string) (locker.Locker, error) {
	lock, err := New(client, id)
	if err != nil {
		return nil, err
	}
	if err := lock.Lock(); err != nil {
		return nil, err
	}
	return lock, nil
}

func New(client *v3.Client, id string) (locker.Locker, error) {
	session, err := concurrency.NewSession(client)
	if err != nil {
		return nil, err
	}
	mutex := concurrency.NewMutex(session, id)
	return &etcd3Lock{
		id:      id,
		mutex:   mutex,
		session: session,
	}, nil
}

func (lock *etcd3Lock) Lock() error {
	return lock.mutex.Lock(context.Background())
}

func (lock *etcd3Lock) Unlock() error {
	var e error
	if err := lock.mutex.Unlock(context.Background()); err != nil {
		e = errors.Wrap(err, "release lock")
	}
	if err := lock.session.Close(); err != nil {
		e = errors.Wrap(err, "close session")
	}
	return e
}
