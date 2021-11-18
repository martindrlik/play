package kafka

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/martindrlik/play/kafka/producer"
)

func Notify(broker, topic string, maxProducers int, acquireTimeout time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	pool := producer.NewPool(broker, maxProducers)
	tryNotify := func(rw http.ResponseWriter, r *http.Request) bool {
		ctx, cancel := context.WithTimeout(r.Context(), acquireTimeout)
		defer cancel()
		p, err := pool.Acquire(ctx)
		if err == r.Context().Err() {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return false
		}
		if err != nil {
			log.Printf("unable to acquire kafka producer: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return false
		}
		defer pool.Release(p)

		delivery := make(chan kafka.Event)
		defer close(delivery)

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny},
			Value: []byte(r.URL.String())},
			delivery)

		e := <-delivery
		m := e.(*kafka.Message) // TODO check type assert: m, ok := e.(*kafka.Message)...

		if m.TopicPartition.Error != nil {
			log.Printf("unable to produce: %v", m.TopicPartition.Error)
			rw.WriteHeader(http.StatusInternalServerError)
			return false
		}

		return true
	}
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			hf(rw, r)
			tryNotify(rw, r)
		}
	}
}
