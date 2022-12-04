package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/martindrlik/play/config"
	"github.com/martindrlik/play/id"
	"github.com/martindrlik/play/metrics"
)

var RequestApiKeyName = "X-Request-ApiKeyName"

// Auth wraps http handler in order to authenticate request.
func Auth(config config.Config, hf http.HandlerFunc) http.HandlerFunc {
	apiKeyNameByValue := make(map[string]string)
	for _, apiKey := range config.ApiKeys {
		apiKeyNameByValue[apiKey.Value] = apiKey.Name
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		apiKeyValue, ok := getRequestApiKey(r)
		apiKeyName := ""
		if ok {
			apiKeyName, ok = apiKeyNameByValue[apiKeyValue]
		}
		if !ok {
			metrics.AuthError()
			if apiKeyValue == "" {
				log.Printf("(%v) missing Authorization header or no value", id.Get(rw))
			} else {
				log.Printf("(%v) invalid api key", id.Get(rw))
			}
			http.Error(rw, "", http.StatusUnauthorized)
			return
		}
		rw.Header().Add(RequestApiKeyName, apiKeyName)
		hf(rw, r)
	}
}

func getRequestApiKey(r *http.Request) (string, bool) {
	v := r.Header.Get("Authorization")
	if !strings.HasPrefix(v, "Bearer ") {
		return "", false
	}
	v = v[len("Bearer "):]
	return v, v != ""
}

// ApiKeyName returns api key name added by Auth.
func ApiKeyName(rw http.ResponseWriter) string {
	return rw.Header().Get(RequestApiKeyName)
}
