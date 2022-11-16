// Play is a http server that enables creating API by uploading go file.
// Go file is then built as a plugin ready to handle requests on
// specified endpoint.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/martindrlik/play/her"
	"github.com/martindrlik/play/kafka"
	"github.com/martindrlik/play/limit"
	"github.com/martindrlik/play/measure"
	"github.com/martindrlik/play/metrics"
	"github.com/martindrlik/play/options"
	"github.com/martindrlik/play/plugin"
	"github.com/martindrlik/play/sequence"
)

var (
	addr                = flag.String("addr", ":8085", "listens on the TCP network address addr")
	maxInFlightRequests = flag.Int("max-in-flight-requests", 250, "limits number of in-flight requests")
	optFile             = flag.String("options", "options.json", "")
	workingDirectory    = flag.String("working-directory", "./wd", "")

	opt options.Options
)

func main() {
	flag.Parse()
	opt = her.Must(options.Load(*optFile))

	if err := os.Chdir(*workingDirectory); err != nil {
		panic(err)
	}

	http.Handle("/metrics", metrics.Handler)
	http.HandleFunc("/upload/", cm(plugin.Upload))
	if opt.Producer.Broker != "" {
		http.HandleFunc("/notify/", cmp(plugin.Run))
	}
	if opt.Consumer.Broker != "" {
		http.HandleFunc("/subscribe/", cmc(plugin.Run))
	}
	http.HandleFunc("/", cm(plugin.Run))

	log.Fatal(http.ListenAndServe(*addr, nil))
}

// cm applies capacity limit and measures duration of hf handler.
func cm(hf http.HandlerFunc) http.HandlerFunc {
	return sequence.Sequence()(limit.Capacity(*maxInFlightRequests)(measure.Measure(hf)))
}

func cmp(hf http.HandlerFunc) http.HandlerFunc {
	return cm(kafka.Notify(opt.Producer)(hf))
}

func cmc(hf http.HandlerFunc) http.HandlerFunc {
	return cm(kafka.Subscribe(context.TODO(), opt.Consumer)(hf))
}
