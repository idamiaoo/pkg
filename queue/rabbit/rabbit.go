package rabbit

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/lunarhalos/pkg/queue"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
}

type rabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	topics     atomic.Value
}

func (rabbit *rabbitMQ) Publish(topic string, message []byte, opts ...queue.Option) error {
	if topic == "" {
		return errors.New("publish: empty topic")
	}
	/*
		if topics, ok := rabbit.topics.Load().(map[string]struct{}); ok {
			if _, ok := topics[topic]; !ok {
				err := rabbit.channel.ExchangeDeclare(topic, amqp.ExchangeTopic, true, false, false, false, nil)
				if err != nil {
					return errors.WithMessage(err, "publish: amqp declare exchange")
				}
				topics[topic] = struct{}{}
				rabbit.topics.Store(topics)
			}
		}
	*/
	q, err := rabbit.channel.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = rabbit.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		},
	)
	if err != nil {
		return errors.WithMessage(err, "publish: amqp publish")
	}
	fmt.Println("publish success")
	return nil
}

func (rabbit *rabbitMQ) Subscribe(topic, group string, handler func(message []byte) error, opts ...queue.Option) error {
	/*
		if topic == "" {
			return errors.New("subscribe: empty topic")
		}

		if topics, ok := rabbit.topics.Load().(map[string]struct{}); ok {
			if _, ok := topics[topic]; !ok {
				err := rabbit.channel.ExchangeDeclare(topic, amqp.ExchangeTopic, true, false, false, false, nil)
				if err != nil {
					return errors.WithMessage(err, "subscribe: amqp declare exchange")
				}
				fmt.Println("ExchangeDeclare")
			}
		}

		q, err := rabbit.channel.QueueDeclare("", false, false, true, false, nil)
		if err != nil {
			return errors.WithMessage(err, "subscribe: amqp declare queue")
		}

		if err := rabbit.channel.QueueBind(q.Name, topic, topic, false, nil); err != nil {
			return errors.WithMessage(err, "subscribe: amqp queue bind")
		}

		msgChan, err := rabbit.channel.Consume(q.Name, group, true, false, false, false, nil)
		if err != nil {
			return errors.WithMessage(err, "subscribe: amqp consume")
		}
		go func() {
			for {
				fmt.Println("等待消息")
				select {
				case message, ok := <-msgChan:
					if !ok {
						return
					}
					fmt.Println(message)
					handler(message.Body)
				}
			}
		}()
		return nil
	*/
	q, err := rabbit.channel.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		nil,
	)
	rabbit.channel.Qos(1, 0, false)
	msgChan, err := rabbit.channel.Consume(q.Name, group, true, false, false, false, nil)
	if err != nil {
		return errors.WithMessage(err, "subscribe: amqp consume")
	}
	go func() {
		for {
			fmt.Println("等待消息")
			select {
			case message, ok := <-msgChan:
				if !ok {
					return
				}
				fmt.Println(message)
				handler(message.Body)
			}
		}
	}()
	return nil
}

func New(c *Config) queue.Queue {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	topics := atomic.Value{}
	topics.Store(map[string]struct{}{})
	return &rabbitMQ{
		connection: conn,
		channel:    ch,
		topics:     topics,
	}
}
