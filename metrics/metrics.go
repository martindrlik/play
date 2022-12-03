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
	unableToCreateId = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_unable_to_create_id_error_count",
		Help: "Unable to create id for request.",
	})
)

// ObserveDuration adds observation to the play_http_request_duration_seconds
// histogram metric.
func ObserveDuration(f float64) { requestDuration.Observe(f) }

// UnableToCreateId increases the unable to create id error count.
func UnableToCreateId() { unableToCreateId.Inc() }
