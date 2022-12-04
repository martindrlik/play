package auth

import (
	"net/http"
	"strings"

	"github.com/martindrlik/play/config"
)

var RequestApiKeyName = "X-Request-ApiKeyName"

// Auth wraps http handler in order to authenticate request.
func Auth(config config.Config, hf http.HandlerFunc) http.HandlerFunc {
	apiKeyNameByValue := make(map[string]string)
	for _, apiKey := range config.ApiKeys {
		apiKeyNameByValue[apiKey.Value] = apiKey.Name
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		apiKey, ok := getRequestApiKey(r)
		if ok {
			apiKey, ok = apiKeyNameByValue[apiKey]
		}
		if !ok {
			// todo metrics
			http.Error(rw, "", http.StatusUnauthorized)
			return
		}
		// metrics
		rw.Header().Add(RequestApiKeyName, apiKey)
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
