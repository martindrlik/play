package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Message struct {
	Value []byte
	Key   []byte
}

func Consume(
	ctx context.Context,
	topic, broker string,
	ch chan<- Message) error {
	kc, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return fmt.Errorf("unable to create kafka consumer: %w", err)
	}
	err = kc.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return fmt.Errorf("unable to subscribe topic: %w", err)
	}
	defer kc.Close()
	for {
		msg, err := kc.ReadMessage(time.Second)
		if err == nil {
			ch <- Message{Value: msg.Value, Key: msg.Key}
		} else if err.(kafka.Error).Code() != kafka.ErrTimedOut {
			// metrics
		}
	}
}
