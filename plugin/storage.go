package plugin

import (
	"net/http"
	"sync"
)

var (
	plugins = map[string]func(http.ResponseWriter, *http.Request){}
	analyze = map[string]error{}

	storageMutex = sync.Mutex{}
)
