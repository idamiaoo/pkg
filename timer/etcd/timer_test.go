package etcd

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/katakurin/pkg/timer"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	v3 "go.etcd.io/etcd/client/v3"
)

func TestTimer(t *testing.T) {
	t1, err := New(v3.Config{
		Endpoints: []string{"192.168.11.248:2379"},
	})
	require.Nil(t, err)

	t2, err := New(v3.Config{
		Endpoints: []string{"192.168.11.248:2379"},
	})
	taskName := "testlf"

	f := func(task []byte) error {
		fmt.Println(string(task))
		return nil
	}

	go func() {
		for {
			<-time.After(time.Second * time.Duration(rand.Int63n(3)+1))
			t1.AddTask(context.TODO(), 4*time.Second, &timer.Task{
				Name:  taskName,
				ID:    ksuid.New().String(),
				Value: []byte("hello1"),
			})
		}
	}()

	go func() {
		for {
			<-time.After(time.Second * time.Duration(rand.Int63n(3)+1))
			t2.AddTask(context.TODO(), 4*time.Second, &timer.Task{
				Name:  taskName,
				ID:    ksuid.New().String(),
				Value: []byte("hello2"),
			})
		}
	}()

	go func() {
		if err := t1.WatchTask(context.TODO(), taskName, f); err != nil {
			log.Println("t1", err)
		}
	}()

	go func() {
		if err := t2.WatchTask(context.TODO(), taskName, f); err != nil {
			fmt.Println("t2", err)
		}
	}()

	select {}
}
