package plugin

import (
	"log"
	"net/http"
	"plugin"
)

func tryLookup(rw http.ResponseWriter, soFile string) (main func(http.ResponseWriter, *http.Request), ok bool) {
	p, err := plugin.Open(soFile)
	if err != nil {
		log.Printf("unable to open plugin %q: %v", soFile, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}
	sym, err := p.Lookup("Main")
	if err != nil {
		log.Printf("unable to lookup Main: %v", err)
		http.Error(rw, "unable to lookup Main", http.StatusBadRequest)
		return nil, false
	}
	main, ok = sym.(func(http.ResponseWriter, *http.Request))
	if !ok {
		log.Printf("unable to type assert to func(http.ResponseWriter, *http.Request): %v", err)
		http.Error(rw, "unable to type assert to func(http.ResponseWriter, *http.Request)",
			http.StatusBadRequest)
		return nil, false
	}
	return
}
