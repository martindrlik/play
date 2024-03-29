package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Handler = promhttp.Handler()

	authError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_auth_error_count",
		Help: "Authentication error.",
	})
	requestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "play_http_request_duration_seconds",
		Help: "Request duration histogram.",
	})
	unableToCreateId = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_unable_to_create_id_error_count",
		Help: "Unable to create id for request.",
	})
	uploadMaxFileLengthExceeded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_upload_max_file_length_exceeded_count",
		Help: "Max upload file length exceeded. Client got status 413.",
	})
	uploadReadingBodyError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_upload_reading_body_error_count",
		Help: "Unable to read request body. Client got status 500 (can be 400).",
	})
	uploadStoringError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_upload_storing_error_count",
		Help: "Unable to store uploaded content. Client got status 500.",
	})
	pluginAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_plugin_added_count",
		Help: "Number of successfully added plugins.",
	})
	pluginError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "play_plugin_error_count",
		Help: "Plugin related (build, lookup, etc.) error count. Client should consult /analyze.",
	})
)

// AuthError increases authentication error count.
func AuthError() { authError.Inc() }

// ObserveDuration adds observation to the play_http_request_duration_seconds
// histogram metric.
func ObserveDuration(f float64) { requestDuration.Observe(f) }

// UnableToCreateId increases the unable to create id error count.
func UnableToCreateId() { unableToCreateId.Inc() }

// UploadMaxFileLengthExceeded increases the upload max file length exceeded count.
func UploadMaxFileLengthExceeded() { uploadMaxFileLengthExceeded.Inc() }

// UploadReadingBodyError increases the upload reading body error count.
func UploadReadingBodyError() { uploadReadingBodyError.Inc() }

// UploadStoringError increases the upload storing error count.
func UploadStoringError() { uploadStoringError.Inc() }

// PluginAdded increases the number of successfully added plugins.
func PluginAdded() { pluginAdded.Inc() }

// PluginError increases the plugin related error count.
func PluginError() { pluginError.Inc() }
