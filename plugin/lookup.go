package plugin

import (
	"fmt"
	"net/http"
	"plugin"
)

func tryLookupHandler(soFile string) (main func(http.ResponseWriter, *http.Request), err error) {
	p, err := plugin.Open(soFile)
	if err != nil {
		err = fmt.Errorf("unable to open plugin %q: %w", soFile, err)
		return
	}
	sym, err := p.Lookup("Main")
	if err != nil {
		err = fmt.Errorf("unable to lookup Main func in %q: %w", soFile, err)
		return
	}
	main, ok := sym.(func(http.ResponseWriter, *http.Request))
	if !ok {
		err = fmt.Errorf("unable to type assert to func(http.ResponseWriter, *http.Request) in %q: %w", soFile, err)
	}
	return
}
