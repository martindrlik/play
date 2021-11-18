package producer

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Pool struct {
	broker string
	sem    chan struct{}
	idle   chan *kafka.Producer
}

func NewPool(broker string, maxProducers int) *Pool {
	sem := make(chan struct{}, maxProducers)
	idle := make(chan *kafka.Producer, maxProducers)
	return &Pool{broker, sem, idle}
}

func (o *Pool) Release(p *kafka.Producer) {
	o.idle <- p
}

func (o *Pool) Acquire(ctx context.Context) (*kafka.Producer, error) {
	select {
	case p := <-o.idle:
		return p, nil
	case o.sem <- struct{}{}:
		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": o.broker})
		if err != nil {
			<-o.sem
		}
		return p, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (o *Pool) TryAcquire(ctx context.Context) (*kafka.Producer, error) {
	for n := 0; n <= 5; n++ {
		p, err := o.Acquire(ctx)
		if err == nil {
			return p, nil
		}
		if err == ctx.Err() {
			return nil, err
		}
		time.Sleep(time.Duration(math.Exp(float64(n))) * time.Second)
	}
	return nil, errors.New("all attempts to acquire producer failed")
}
