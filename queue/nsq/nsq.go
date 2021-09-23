package nsq

import (
	"log"
	"os"

	"github.com/lunarhalos/pkg/queue"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type nsqMQ struct {
	producer *nsq.Producer
	config   *Config
}

type Config struct {
	Nsqd       string
	Lookupd    string
	AuthSecret string
}

func (mq *nsqMQ) Publish(topic string, message []byte, opts ...queue.Option) error {
	return mq.producer.Publish(topic, message)
}

type nsqHandler struct {
	process func(message []byte) error
}

func (h *nsqHandler) HandleMessage(m *nsq.Message) error {
	return h.process(m.Body)
}

func (mq *nsqMQ) Subscribe(topic, group string, handler func(message []byte) error, opts ...queue.Option) error {
	config := nsq.NewConfig()
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
	config.Set("auth_secret", c.AuthSecret)
	producer, err := nsq.NewProducer(c.Nsqd, config)
	if err != nil {
		logrus.Fatal(err)
	}
	producer.SetLogger(log.New(os.Stdout, "", log.Flags()), nsq.LogLevelInfo)

	return &nsqMQ{
		config:   c,
		producer: producer,
	}
}
