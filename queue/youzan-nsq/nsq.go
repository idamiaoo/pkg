package nsq

import (
	"fmt"
	"log"
	"os"

	"github.com/pescaria/pkg/queue"

	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
	"github.com/youzan/go-nsq"
)

type nsqMQ struct {
	producer *nsq.TopicProducerMgr
	config   *Config
}

type Config struct {
	Nsqd       string
	Lookupd    string
	AuthSecret string
}

func (mq *nsqMQ) Publish(topic string, message []byte, opts ...queue.Option) (err error) {
	baseConfig := &queue.Config{}
	for _, opt := range opts {
		opt(baseConfig)
	}
	if baseConfig.OrderID != nil {
		_, _, _, err = mq.producer.PublishOrdered(topic, baseConfig.OrderID, message)
	} else {
		err = mq.producer.Publish(topic, message)
	}
	return err
}

type nsqHandler struct {
	process func(message []byte) error
}

func (h *nsqHandler) HandleMessage(m *nsq.Message) error {
	return h.process(m.Body)
}

func (mq *nsqMQ) Subscribe(topic, group string, handler func(message []byte) error, opts ...queue.Option) error {
	baseConfig := &queue.Config{}
	for _, opt := range opts {
		opt(baseConfig)
	}
	fmt.Printf("baseConfig %v", *baseConfig)
	config := nsq.NewConfig()
	if baseConfig.EnableOrdered {
		config.EnableOrdered = true
		config.Hasher = murmur3.New32()
	}
	consumer, err := nsq.NewConsumer(topic, group, config)
	if err != nil {
		return err
	}
	consumer.SetLogger(log.New(os.Stdout, "", log.Flags()), nsq.LogLevelInfo)
	consumer.AddHandler(&nsqHandler{process: handler})
	err = consumer.ConnectToNSQLookupd(mq.config.Lookupd)
	if err != nil {
		return err
	}
	return nil
}

func New(c *Config) queue.Queue {
	config := nsq.NewConfig()
	config.EnableOrdered = true
	config.Hasher = murmur3.New32()
	config.Set("auth_secret", c.AuthSecret)
	producer, err := nsq.NewTopicProducerMgr([]string{}, config)
	if err != nil {
		logrus.Fatal(err)
	}

	producer.SetLogger(log.New(os.Stdout, "", log.Flags()), nsq.LogLevelInfo)
	producer.ConnectToNSQLookupd(c.Lookupd)
	return &nsqMQ{
		config:   c,
		producer: producer,
	}
}
