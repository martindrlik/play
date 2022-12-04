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
	// api key authenticated
	h.HandleFunc("/upload/", id.Gen(acm(config, plugin.Upload(produce))))
	h.HandleFunc("/analyze/", id.Gen(acm(config, plugin.Analyze)))
	h.HandleFunc("/", id.Gen(acm(config, plugin.Execute)))
	return h
}

// acm adds request authentication to hf otherwise the same as cm.
func acm(config config.Config, hf http.HandlerFunc) http.HandlerFunc {
	return auth.Auth(config, cm(config, hf))
}

// cm wraps hf to limit in-fligh requests and to measure handler's performance.
func cm(config config.Config, hf http.HandlerFunc) http.HandlerFunc {
	return limit.Capacity(config.RequestLimit)( // 429
		measure.Measure( // logs duration
			metrics.ObserveDuration,
			hf))
}
