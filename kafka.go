package main

import (
	"context"
	"fmt"

	"github.com/martindrlik/play/config"
	"github.com/martindrlik/play/her"
	"github.com/martindrlik/play/kafka"
	"github.com/martindrlik/play/plugin"
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
			case m := <-m:
				plugin.Consume(m.Value, string(m.Key))
			case <-ctx.Done():
				panic(fmt.Errorf("canceled: %w", ctx.Err()))
			}
		}
	}()
	her.Must1(kafka.Consume(
		ctx,
		config.KafkaUploadTopic,
		config.KafkaBroker,
		m))
}
