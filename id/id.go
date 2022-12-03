package id

import (
	"context"
	"log"
	"net/http"

	"github.com/martindrlik/play/metrics"
	"github.com/segmentio/ksuid"
)

var RequestIdHeaderName = "X-Request-Id"

// Gen creates unique id and adds it to response headers
// under header name given by RequestIdHeaderName.
func Gen() func(http.HandlerFunc) http.HandlerFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {
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
}

// Get returns id from context added by Gen.
func Get(ctx context.Context) string {
	if rid, ok := ctx.Value(RequestIdHeaderName).(string); ok {
		return rid
	}
	return ""
}
