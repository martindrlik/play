package id

import (
	"log"
	"net/http"

	"github.com/martindrlik/play/metrics"
	"github.com/segmentio/ksuid"
)

var RequestIdHeaderName = "X-Request-Id"

// Gen wraps hf in order to creates unique id and adds
// it to response headers under header name given by
// RequestIdHeaderName.
func Gen(hf http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, err := ksuid.NewRandom()
		if err == nil {
			rw.Header().Add(RequestIdHeaderName, id.String())
			hf(rw, r)
		} else {
			metrics.UnableToCreateId()
			log.Printf("unable to create unique id: %v", err)
			hf(rw, r)
		}
	}
}

// Get returns id added by Gen.
func Get(rw http.ResponseWriter) string {
	return rw.Header().Get(RequestIdHeaderName)
}
