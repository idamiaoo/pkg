package etcdlocker

import (
	"testing"
	"time"

	v3 "go.etcd.io/etcd/client/v3"

	"github.com/stretchr/testify/assert"
)

func TestEtcd3Locker(t *testing.T) {
	a := assert.New(t)

	client, err := v3.New(v3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Unable to connect to etcd3: %v", err)
	}
	defer client.Close()

	id := "abcdefgh"

	lock1, err := New(client, id)
	a.NoError(err)
	lock1.Lock()
	lock2, err := New(client, id)
	a.NoError(err)

	t.Log("lock1 acquire lock")
	time.AfterFunc(10*time.Second, func() {
		lock1.Unlock()
		t.Log("lock1 release lock")
	})
	lock1.Lock()
	t.Log("lock1 re acquire lock")
	lock2.Lock()
	t.Log("lock2 acquire lock")
	lock2.Unlock()
	t.Log("lock2 release lock")
}
