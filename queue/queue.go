package queue

import (
	"fmt"
	"hash"
)

type Queue interface {
	Publish(topic string, message []byte, opts ...Option) error
	Subscribe(topic, group string, handler func(message []byte) error, opts ...Option) error
}

type Config struct {
	EnableTrace   bool
	EnableOrdered bool
	TraceID       uint64
	OrderID       []byte
	Hasher        hash.Hash32
}

type Option func(c *Config)

func WithOrderedID(orderID []byte) Option {
	return func(c *Config) {
		c.OrderID = orderID
	}
}

func EnableTrace() Option {
	return func(c *Config) {
		c.EnableTrace = true
	}
}

func EnableOrdered() Option {
	return func(c *Config) {
		fmt.Println("222222222222222222222222222")
		c.EnableOrdered = true
	}
}

func WithHasher(hasher hash.Hash32) Option {
	return func(c *Config) {
		c.Hasher = hasher
	}
}
