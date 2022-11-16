package kafka

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/martindrlik/play/backoff"
	"github.com/martindrlik/play/options"
)

func Subscribe(ctx context.Context, consumerOpt options.KafkaOptions) func(http.HandlerFunc) http.HandlerFunc {
	consPool := sync.Pool{
		New: func() any {
			c, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": consumerOpt.Broker})
			if err != nil {
				log.Printf("unable to create kafka consumer: %v", err)
				return nil
			}
			return c
		},
	}
	return func(hf http.HandlerFunc) http.HandlerFunc {
		go pull(ctx, &consPool, consumerOpt.Topic, hf)
		return func(rw http.ResponseWriter, r *http.Request) {
			hf(rw, r)
		}
	}
}

func tryPull(ctx context.Context, pool *sync.Pool, topic string, hf http.HandlerFunc) bool {
	c := pool.Get().(*kafka.Consumer)
	if c == nil {
		return false
	}
	// TODO in error case we probably need hijack "broken consumer"
	defer pool.Put(c)
	err := c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Printf("unable to subscribe topic: %v", err)
		return false
	}
	defer c.Unsubscribe()
	commit := func() {
		_, err := c.Commit()
		if err != nil {
			log.Printf("unable to commit offsets: %v", err)
		}
	}

	messages := make([]string, 0, 5)
	var cancel context.CancelFunc
	for {
		select {
		case <-ctx.Done():
			commit()
			goto proc
		default:
			switch x := c.Poll(100).(type) {
			case nil:
			case *kafka.Message:
				from := string(x.Headers[0].Value)
				text := string(x.Value)
				messages = append(messages, fmt.Sprintf("message: %s, from: %s", text, from))
				if len(messages) == cap(messages) {
					commit()
					goto proc
				}
				if len(messages) == 1 {
					linger := 5 * time.Millisecond
					log.Printf("pulled one message, waiting %v for more", linger)
					ctx, cancel = context.WithTimeout(ctx, linger)
					defer cancel()
				}
			default:
				log.Printf("pulled something else: %v", x)
			}
		}
	}
proc:
	log.Printf("pulled %v messages, TODO implement hf(rw, r) for each", len(messages))
	return true
}

func pull(ctx context.Context, pool *sync.Pool, topic string, hf http.HandlerFunc) {
	retry := 0
	for {
		if !tryPull(ctx, pool, topic, hf) {
			d := backoff.Exp(retry)
			log.Printf("unable to tryPull: another attempt after %v", d)
			time.Sleep(d)
			if d < 18*time.Minute {
				retry++
			}
			continue
		}
		retry = 0
	}
}
