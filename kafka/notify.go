package kafka

import (
	"log"
	"net/http"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/martindrlik/play/options"
)

// Notify
func Notify(producerOpt options.KafkaOptions) func(http.HandlerFunc) http.HandlerFunc {
	prodPool := sync.Pool{
		New: func() any {
			p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": producerOpt.Broker})
			if err != nil {
				log.Printf("unable to create kafka producer: %v", err)
				return nil
			}
			return p
		},
	}
	tryNotify := func(rw http.ResponseWriter, r *http.Request) bool {
		p := prodPool.Get().(*kafka.Producer)
		if p == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return false
		}
		defer prodPool.Put(p)

		delivery := make(chan kafka.Event)
		defer close(delivery)

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &producerOpt.Topic,
				Partition: kafka.PartitionAny},
			Value: []byte(r.URL.String())},
			delivery)
		if err != nil {
			log.Printf("unable to produce message: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return false
		}

		e := <-delivery
		m := e.(*kafka.Message) // TODO check type assert: m, ok := e.(*kafka.Message)...

		if m.TopicPartition.Error != nil {
			log.Printf("topic partition error: %v", m.TopicPartition.Error)
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
