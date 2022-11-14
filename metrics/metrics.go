package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Handler = promhttp.Handler()

	requestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "play_http_request_duration_seconds",
		Help: "Request duration histogram.",
	})
)

// ObserveDuration adds observation to the play_http_request_duration_seconds
// histogram metric.
func ObserveDuration(f float64) { requestDuration.Observe(f) }
