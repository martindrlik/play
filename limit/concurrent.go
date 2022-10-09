package limit

import "net/http"

// Concurrent limits number of request processed at the same time.
func Concurrent(max int) func(http.HandlerFunc) http.HandlerFunc {
	ch := make(chan struct{}, max)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			select {
			case ch <- struct{}{}:
				defer func() { <-ch }()
			default:
				rw.WriteHeader(http.StatusTooManyRequests)
				return
			}
			hf(rw, r)
		}
	}
}
