package main

import (
	"context"

	"github.com/martindrlik/play/config"
	"github.com/martindrlik/play/her"
	"github.com/martindrlik/play/kafka"
)

func producer(ctx context.Context, config config.Config) func(value, key []byte) error {
	p := her.Must(kafka.NewProducer(ctx, config.KafkaBroker))
	return func(value, key []byte) error {
		return p.Produce(config.KafkaUploadTopic, value, key)
	}
}

func consumeMessages(ctx context.Context, config config.Config) {
	m := make(chan kafka.Message)
	go func() {
		for {
			select {
			case <-m:
				panic("not implemented")
			case <-ctx.Done():
				panic("canceled")
			}
		}
	}()
	her.Must1(kafka.Consume(
		ctx,
		config.KafkaUploadTopic,
		config.KafkaBroker,
		m))
}
