package measure

import (
	"log"
	"net/http"
	"time"

	"github.com/martindrlik/play/id"
)

// Measure wraps http handler in order to measure and logs its duration.
func Measure(observeDuration func(f float64), hf http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		hf(rw, r)
		elapsed := time.Since(start)
		observeDuration(elapsed.Seconds())
		log.Printf("(%v) %s%s took %v", id.Get(rw), r.Host, r.URL, elapsed)
	}
}
