package limit

import "net/http"

// Capacity wraps http handler in order to limit maximum
// number of in-flight request to max.
func Capacity(max uint) func(http.HandlerFunc) http.HandlerFunc {
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
