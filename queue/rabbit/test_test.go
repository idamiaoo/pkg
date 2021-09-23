package rabbit

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/lunarhalos/pkg/queue"

	"github.com/stretchr/testify/require"
)

var c Config = Config{
	Host:     "192.168.11.248",
	Port:     5672,
	Username: "guest",
	Password: "guest",
}

func TestRabbit(t *testing.T) {
	q1 := New(&c)
	q2 := New(&c)
	q3 := New(&c)
	q4 := New(&c)
	var err error
	err = q4.Publish("test.beetle.aa1aaa", []byte("hello"))
	require.Nil(t, err)
	err = q1.Subscribe("test.beetle.aa1aaa", "a", func(message []byte) error {
		fmt.Println("q1: ", string(message))
		return nil
	})
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

func TestCourse(t *testing.T) {
	q1 := New(&c)
	var err error
	err = q1.Subscribe(queue.TopicVenusCoursePublish, "", func(message []byte) error {
		fmt.Println("q1: ", string(message))
		return nil
	})
	require.Nil(t, err)
	select {}

}

func TestAscII(t *testing.T) {
	st := []string{
		"海底泡泡",
		"大转盘",
		"板下活动",
		"语音小喇叭",
		"人脸识别抢星星",
		"多场景AR",
		"语音选择",
		"泡泡游戏深海版",
		"泡泡游戏太空版",
		"泡泡游戏农场版",
		"语音抢热气球山川版",
		"语音抢热气球白云版",
		"语音抢热气球海洋",
		"markstarcontent",
		"kc星空ho'w",
		"分组问答",
	}
	sort.Strings(st)
	for _, v := range st {
		fmt.Println(v)
	}
}
