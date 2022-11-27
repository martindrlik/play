package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	kp *kafka.Producer
}

// NewProducer creates new Producer.
func NewProducer(ctx context.Context, broker string) (*Producer, error) {
	kp, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, fmt.Errorf("unable to create producer: %w", err)
	}
	p := &Producer{kp}
	go p.pullEvents(ctx)
	return p, nil
}

func (p *Producer) pullEvents(ctx context.Context) {
	for {
		select {
		case event := <-p.kp.Events():
			switch ev := event.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("failed to produce message: %v", ev.TopicPartition)
				}
			}
		case <-ctx.Done():
			log.Printf("pulling producer events canceled: %v", ctx.Err())
		}
	}
}

// Produce produces single message given by value and key to topic.
func (p *Producer) Produce(topic string, value, key []byte) error {
	err := p.kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny},
		Value: value,
		Key:   key},
		nil)
	if err != nil {
		return fmt.Errorf("unable to produce message: %w", err)
	}
	return nil
}
