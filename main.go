// Play is a http server that enables creating API by uploading
// handler's source code. Compiled as a plugin ready to handle
// requests on specified endpoint.
package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/martindrlik/play/auth"
	"github.com/martindrlik/play/config"
	"github.com/martindrlik/play/her"
	"github.com/martindrlik/play/id"
	"github.com/martindrlik/play/limit"
	"github.com/martindrlik/play/measure"
	"github.com/martindrlik/play/metrics"
	"github.com/martindrlik/play/plugin"
)

var (
	addr = flag.String("addr", ":8085", "listens on the TCP network address addr")
	conf = flag.String("config", "config.json", "configuration file")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := her.Must(os.Open(*conf))
	defer cfg.Close()

	config := her.Must(config.Load(cfg))
	go consumeMessages(ctx, config)
	produce := producer(ctx, config)
	her.Must1(http.ListenAndServe(*addr, handler(config, produce)))
}

func handler(config config.Config, produce func(value, key []byte) error) http.Handler {
	h := http.NewServeMux()
	h.Handle("/metrics", metrics.Handler)
	//
	h.HandleFunc("/upload/", cm(config, plugin.Upload(produce)))
	h.HandleFunc("/analyze/", cm(config, plugin.Analyze))
	h.HandleFunc("/", cm(config, plugin.Execute))
	return h
}

func cm(config config.Config, hf http.HandlerFunc) http.HandlerFunc {
	return id.Gen()( // add X-Request-Id
		auth.Auth( // 401 or add X-Request-ApiKeyName
			config,
			limit.Capacity(config.RequestLimit)( // 429 if no room
				measure.Measure( // logs duration
					metrics.ObserveDuration,
					hf))))
}
