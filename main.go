// Play is a http server that enables creating API by uploading
// handler's source code. Compiled as a plugin ready to handle
// requests on specified endpoint.
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
	addr = flag.String("addr", ":8085", "listens on the TCP network address addr")
	rcap = flag.Int("limit", 250, "limits number of in-flight requests to limit")
	opts = flag.String("options", "options.json", "")
	wdir = flag.String("wd", "", "working directory")
)

func main() {
	flag.Parse()
	if *wdir != "" {
		her.Must1(os.Chdir(*wdir))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	o := her.Must(options.Load(*opts))
	go consume(ctx, o)

	http.Handle("/metrics", metrics.Handler)
	http.HandleFunc("/upload/", cm(plugin.Upload(ctx, o)))
	http.HandleFunc("/", cm(plugin.Run))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func consume(ctx context.Context, o options.Options) {
	m := make(chan kafka.Message)
	go func() {
		for {
			select {
			case <-m:
				panic("not implemented")
			case <-ctx.Done():
				panic("canceled")
			}
		}
	}()
	her.Must1(kafka.Consume(
		ctx,
		o.KafkaUploadTopic,
		o.KafkaBroker,
		m))
}

func cm(hf http.HandlerFunc) http.HandlerFunc {
	return sequence.Sequence()(limit.Capacity(*rcap)(measure.Measure(hf)))
}
