package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martindrlik/play/auth"
	"github.com/martindrlik/play/config"
)

func TestAuth(t *testing.T) {
	mainApiKey := config.ApiKey{Name: "main-api-key", Value: "secret"}
	config := config.Config{
		ApiKeys: []config.ApiKey{mainApiKey},
	}
	t.Run("no authorization header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		auth.Auth(config, func(rw http.ResponseWriter, r *http.Request) {
			panic("should not be called")
		})(rec, req)
		actual := rec.Result().StatusCode
		if actual != http.StatusUnauthorized {
			t.Errorf("expected status code %v got %v", http.StatusUnauthorized, actual)
		}
	})
	t.Run("invalid api key value", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add("Authorization", "Bearer invalid"+mainApiKey.Value)
		auth.Auth(config, func(rw http.ResponseWriter, r *http.Request) {
			panic("should not be called")
		})(rec, req)
		actual := rec.Result().StatusCode
		if actual != http.StatusUnauthorized {
			t.Errorf("expected status code %v got %v", http.StatusUnauthorized, actual)
		}
	})
	t.Run("valid api key", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add("Authorization", "Bearer "+mainApiKey.Value)
		called := 0
		wantStatus := http.StatusFound
		auth.Auth(config, func(rw http.ResponseWriter, r *http.Request) {
			called++
			rw.WriteHeader(wantStatus)
		})(rec, req)
		if called != 1 {
			t.Errorf("expected handler to be called once, called %v", called)
		}
		if actual := rec.Result().StatusCode; actual != wantStatus {
			t.Errorf("expected status code %v got %v", wantStatus, actual)
		}
		if actual := rec.Header().Get(auth.RequestApiKeyName); actual != mainApiKey.Name {
			t.Errorf("expected api key name header to be %q, got %q", mainApiKey.Name, actual)
		}
	})
}
