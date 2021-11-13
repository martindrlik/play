package measure

import (
	"log"
	"net/http"
	"time"

	"github.com/martindrlik/play/sequence"
)

func Measure(hf http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		hf(rw, r)
		log.Printf("(%d) %s%s took %v", sequence.Get(r.Context()), r.Host, r.URL, time.Now().Sub(start))
	}
}
