package plugin

import (
	"fmt"
	"net/http"
)

// Execute looks up handler by r.URL.Path and calls the handler.
func Execute(rw http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	main, ok := func() (main func(http.ResponseWriter, *http.Request), ok bool) {
		storageMutex.Lock()
		defer storageMutex.Unlock()
		main, ok = plugins[name]
		return
	}()
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	main(rw, r)
}

// Analyze looks up error for API given by r.URL.Path and responds
// by error message or 404 if no error found.
func Analyze(rw http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/analyze"):]
	err, ok := func() (err error, ok bool) {
		storageMutex.Lock()
		defer storageMutex.Unlock()
		err, ok = analyze[name]
		return
	}()
	if !ok || err == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprintf(rw, "%v\n", err)
}
