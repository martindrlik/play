package measure

import (
	"log"
	"net/http"
	"time"

	"github.com/martindrlik/play/metrics"
	"github.com/martindrlik/play/sequence"
)

func Measure(hf http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		hf(rw, r)
		elapsed := time.Since(start)
		metrics.RequestDurationObserve(elapsed.Seconds())
		log.Printf("measure: %v %s%s took %v", sequence.Get(r.Context()), r.Host, r.URL, elapsed)
	}
}
