package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martindrlik/play/auth"
)

func TestAuth(t *testing.T) {
	const (
		mainApiKeyValue = "main-api-key-secret"
		mainApiKeyName  = "main-api-key"
	)
	nameByApiKey := map[string]string{
		mainApiKeyValue: mainApiKeyName,
	}
	t.Run("no authorization header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		auth.Auth(nameByApiKey, func(rw http.ResponseWriter, r *http.Request) {
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
		req.Header.Add("Authorization", "Bearer invalid"+mainApiKeyValue)
		auth.Auth(nameByApiKey, func(rw http.ResponseWriter, r *http.Request) {
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
		req.Header.Add("Authorization", "Bearer "+mainApiKeyValue)
		called := 0
		wantStatus := http.StatusFound
		auth.Auth(nameByApiKey, func(rw http.ResponseWriter, r *http.Request) {
			called++
			rw.WriteHeader(wantStatus)
		})(rec, req)
		if called != 1 {
			t.Errorf("expected handler to be called once, called %v", called)
		}
		if actual := rec.Result().StatusCode; actual != wantStatus {
			t.Errorf("expected status code %v got %v", wantStatus, actual)
		}
		if actual := rec.Header().Get(auth.RequestApiKeyName); actual != mainApiKeyName {
			t.Errorf("expected api key name header to be %q, got %q", mainApiKeyName, actual)
		}
		if actual := auth.ApiKeyName(rec); actual != mainApiKeyName {
			t.Errorf("expected api key name returned by auth.ApiKeyName to be %q, got %q", mainApiKeyName, actual)
		}
	})
}
