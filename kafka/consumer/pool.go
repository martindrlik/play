package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/martindrlik/play/backoff"
)

type Pool struct {
	broker string
	sem    chan struct{}
	idle   chan *kafka.Consumer
}

func NewPool(broker string, maxConsumers int) *Pool {
	sem := make(chan struct{}, maxConsumers)
	idle := make(chan *kafka.Consumer, maxConsumers)
	return &Pool{broker, sem, idle}
}

func (o *Pool) Release(c *kafka.Consumer) {
	o.idle <- c
}

func (o *Pool) Acquire(ctx context.Context) (*kafka.Consumer, error) {
	select {
	case c := <-o.idle:
		return c, nil
	case o.sem <- struct{}{}:
		c, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": o.broker})
		if err != nil {
			<-o.sem
		}
		return c, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (o *Pool) TryAcquire(ctx context.Context) (*kafka.Consumer, error) {
	for n := 0; n <= 5; n++ {
		c, err := o.Acquire(ctx)
		if err == nil {
			return c, nil
		}
		if err == ctx.Err() {
			return nil, ctx.Err()
		}
		time.Sleep(backoff.Exp(n))
	}
	return nil, errors.New("all attempts to acquire consumer failed")
}
