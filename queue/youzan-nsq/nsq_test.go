package nsq

import (
	"fmt"
	"testing"
	"time"

	"github.com/lunarhalos/pkg/queue"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

var c Config = Config{
	// Nsqd:    "192.168.11.8:4150",
	// Lookupd: "192.168.11.8:4161",
	// Nsqd:    "127.0.0.1:4150",
	Lookupd: "127.0.0.1:4161",
}

func TestNSQ(t *testing.T) {
	q1 := New(&c)
	q2 := New(&c)
	q3 := New(&c)
	q4 := New(&c)
	var err error
	err = q4.Publish("test.beetle.aa1aaa", []byte("hello"))
	require.Nil(t, err)
	err = q1.Subscribe("test.beetle.aa1aaa", "aa", func(message []byte) error {
		fmt.Println("q1: ", string(message))
		return nil
	})
	if err != nil {
		println("sq:", err)
	}

	require.Nil(t, err)
	err = q2.Subscribe("test.beetle.aa1aaa", "aa", func(message []byte) error {
		fmt.Println("q2: ", string(message))
		return nil
	})
	require.Nil(t, err)
	go func() {
		var i int
		for {
			i++
			fmt.Println("push message")
			msg := fmt.Sprintf("hello %d", i)
			q3.Publish("test.beetle.aa1aaa", []byte(msg))
			time.Sleep(1 * time.Second)
		}
	}()

	select {}

}

func TestNSQOrder(t *testing.T) {
	q1 := New(&c)
	q2 := New(&c)
	q3 := New(&c)
	q4 := New(&c)
	var err error
	err = q4.Publish("test.beetle.aa1aaa.ordered", []byte("hello"), queue.WithOrderedID(ksuid.New().Bytes()))
	require.Nil(t, err)
	err = q1.Subscribe("test.beetle.aa1aaa.ordered", "aa", func(message []byte) error {
		fmt.Println("q1: ", string(message))
		return nil
	}, queue.EnableOrdered())
	if err != nil {
		println("sq:", err)
	}

	require.Nil(t, err)
	err = q2.Subscribe("test.beetle.aa1aaa.ordered", "aa", func(message []byte) error {
		fmt.Println("q2: ", string(message))
		return nil
	}, queue.EnableOrdered())
	require.Nil(t, err)
	go func() {
		var i int
		for {
			i++
			fmt.Println("push message")
			msg := fmt.Sprintf("hello %d", i)
			q3.Publish("test.beetle.aa1aaa.ordered", []byte(msg), queue.WithOrderedID(ksuid.New().Bytes()))
			time.Sleep(1 * time.Second)
		}
	}()

	select {}

}
