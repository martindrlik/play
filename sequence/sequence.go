package sequence

import (
	"context"
	"net/http"
	"sync"
)

func Sequence(start int64) func(http.HandlerFunc) http.HandlerFunc {
	n := start
	m := sync.Mutex{}
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			m.Lock()
			n++
			seq := n
			m.Unlock()
			ctx := context.WithValue(r.Context(), "sequence_int64", seq)
			hf(rw, r.WithContext(ctx))
		}
	}
}

func Get(ctx context.Context) int64 {
	if seq, ok := ctx.Value("sequence_int64").(int64); ok {
		return seq
	}
	return 0
}
