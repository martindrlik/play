package plugin

import (
	"net/http"
	"sync"
)

var (
	plugins      = map[string]func(http.ResponseWriter, *http.Request){}
	pluginsMutex = sync.Mutex{}
)
