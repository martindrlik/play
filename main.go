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
	addr  = flag.String("addr", ":8085", "listens on the TCP network address addr")
	rcap  = flag.Int("capacity", 250, "limits number of in-flight requests")
	ofile = flag.String("options", "options.json", "")
	wdir  = flag.String("working-directory", "./wd", "")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()
	o := her.Must(options.Load(*ofile))

	if err := os.Chdir(*wdir); err != nil {
		panic(err)
	}

	p := her.Must(kafka.NewProducer(ctx, o.KafkaBroker))
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
	go func() {
		err := kafka.Consume(
			ctx,
			o.KafkaUploadTopic,
			o.KafkaBroker,
			m)
		panic(err)
	}()

	http.Handle("/metrics", metrics.Handler)
	http.HandleFunc("/upload/", cm(plugin.Upload(o.KafkaUploadTopic, p)))
	http.HandleFunc("/", cm(plugin.Run))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func cm(hf http.HandlerFunc) http.HandlerFunc {
	return sequence.Sequence()(limit.Capacity(*rcap)(measure.Measure(hf)))
}
