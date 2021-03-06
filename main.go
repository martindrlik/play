package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/martindrlik/play/kafka"
	"github.com/martindrlik/play/limit"
	"github.com/martindrlik/play/measure"
	"github.com/martindrlik/play/options"
	"github.com/martindrlik/play/plugin"
	"github.com/martindrlik/play/sequence"
)

var (
	addr    = flag.String("addr", ":8085", "")
	start   = flag.Int64("start", 1000, "every request is identified by increasing number, start sets initial value")
	max     = flag.Int("max", 10, "max sets limit of how many requests can be processed in one time")
	optfile = flag.String("options", "options.json", "")

	opt options.Options
)

func main() {
	flag.Parse()
	opt = options.Must(options.Load(*optfile))
	http.HandleFunc("/upload/", mc(plugin.Upload))
	http.HandleFunc("/notify/", mnc(plugin.Run))
	http.HandleFunc("/subscribe/", msc(plugin.Run))
	http.HandleFunc("/", mc(plugin.Run))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func mc(hf http.HandlerFunc) http.HandlerFunc {
	return sequence.Sequence(*start)(measure.Measure(limit.Concurrent(*max)(hf)))
}

func mnc(hf http.HandlerFunc) http.HandlerFunc {
	return mc(kafka.Notify(opt.Producer)(hf))
}

func msc(hf http.HandlerFunc) http.HandlerFunc {
	return mc(kafka.Subscribe(context.TODO(), opt.Consumer)(hf))
}
