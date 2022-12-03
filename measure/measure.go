package measure

import (
	"log"
	"net/http"
	"time"

	"github.com/martindrlik/play/id"
	"github.com/martindrlik/play/metrics"
)

// Measure wraps http handler in order to measure and logs its duration.
func Measure(hf http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		hf(rw, r)
		elapsed := time.Since(start)
		metrics.ObserveDuration(elapsed.Seconds())
		log.Printf("(%v) %s%s took %v", id.Get(rw), r.Host, r.URL, elapsed)
	}
}
