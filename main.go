package main

import (
	"log"
	"net/http"

	"github.com/martindrlik/play/limit"
	"github.com/martindrlik/play/measure"
	"github.com/martindrlik/play/public"
	"github.com/martindrlik/play/sequence"
)

func main() {
	http.HandleFunc("/api/0/run", mc5(public.Run))
	log.Fatal(http.ListenAndServe(":8070", nil))
}

func mc5(hf http.HandlerFunc) http.HandlerFunc {
	return sequence.Sequence(1_000)(measure.Measure(limit.Concurrent(5)(hf)))
}
