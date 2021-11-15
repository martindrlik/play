package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/martindrlik/play/limit"
	"github.com/martindrlik/play/measure"
	"github.com/martindrlik/play/plugin"
	"github.com/martindrlik/play/sequence"
)

var (
	addr  = flag.String("addr", ":8085", "")
	start = flag.Int64("start", 1000, "every request is identified by increasing number, start sets initial value")
	max   = flag.Int("max", 10, "max sets limit of how many requests can be processed in one time")
)

func main() {
	http.HandleFunc("/upload/", mc5(plugin.Upload))
	http.HandleFunc("/", mc5(plugin.Run))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func mc5(hf http.HandlerFunc) http.HandlerFunc {
	return sequence.Sequence(*start)(measure.Measure(limit.Concurrent(*max)(hf)))
}
