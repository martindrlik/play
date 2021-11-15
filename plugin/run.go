package plugin

import (
	"net/http"
)

func Run(rw http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	main, ok := func() (main func(http.ResponseWriter, *http.Request), ok bool) {
		pluginsMutex.Lock()
		defer pluginsMutex.Unlock()
		main, ok = plugins[name]
		return
	}()
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	main(rw, r)
}
