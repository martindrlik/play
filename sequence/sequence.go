package sequence

import (
	"context"
	"log"
	"net/http"

	"github.com/segmentio/ksuid"
)

func Sequence() func(http.HandlerFunc) http.HandlerFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			id, err := ksuid.NewRandom()
			if err == nil {
				ctx := context.WithValue(r.Context(), "sequence_id", id.String())
				hf(rw, r.WithContext(ctx))
			} else {
				log.Printf("unable to generate id: %v", err)
				hf(rw, r)
			}
		}
	}
}

func Get(ctx context.Context) string {
	if id, ok := ctx.Value("sequence_id").(string); ok {
		return id
	}
	return ""
}
