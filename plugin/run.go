package plugin

import (
	"log"
	"net/http"
	"path"
)

func Run(rw http.ResponseWriter, r *http.Request) {
	key := path.Base(r.URL.Path)
	log.Print(key)
	main, ok := func() (main func(http.ResponseWriter, *http.Request), ok bool) {
		pluginsMutex.Lock()
		defer pluginsMutex.Unlock()
		main, ok = plugins[key]
		return
	}()
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	main(rw, r)
}
