package limit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martindrlik/play/limit"
)

func TestCapacity(t *testing.T) {
	for _, tc := range []struct {
		name string
		max  uint
	}{
		{"zero request allowed", 0},
		{"one request allowed", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// note that all in-flight requests are waiting for done
			// one request is expected to get 429 too many requests
			// this request does not enter so it is not blocked
			// this request closes done channel, unblocks
			// in-flight requests, check their status,
			// test ends
			done := make(chan struct{})
			statuses := make(chan int)
			h := limit.Capacity(tc.max)(func(w http.ResponseWriter, r *http.Request) {
				<-done
			})
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			for i := uint(0); i < tc.max+1; i++ {
				go func() {
					rec := httptest.NewRecorder()
					h(rec, req)
					if rec.Result().StatusCode == http.StatusTooManyRequests {
						close(done) // unblock
					} else {
						statuses <- rec.Result().StatusCode
					}
				}()
			}
			for i := uint(0); i < tc.max; i++ {
				if actual := <-statuses; actual != http.StatusOK {
					t.Errorf("expected 200 got %v", actual)
				}
			}
		})
	}
}
