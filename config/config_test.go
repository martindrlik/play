package config_test

import (
	"strings"
	"testing"

	"github.com/martindrlik/play/config"
	"github.com/martindrlik/play/her"
)

func TestLoad(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		config := her.Must(config.Load(strings.NewReader(`{}`)))
		if want := uint(250); config.RequestLimit != want {
			t.Errorf("expected default request limit to be %v got %v", want, config.RequestLimit)
		}
	})
	t.Run("custom defaults", func(t *testing.T) {
		config.DefaultRequestLimit = 1000
		config := her.Must(config.Load(strings.NewReader(`{}`)))
		if want := uint(1000); config.RequestLimit != want {
			t.Errorf("expected default request limit to be %v got %v", want, config.RequestLimit)
		}
	})
	t.Run("api keys", func(t *testing.T) {
		config.DefaultRequestLimit = 1000
		config := her.Must(config.Load(strings.NewReader(`{
			"apiKeys": [
				{ "name": "main-api-key", "value": "foo" }
			]
		}`)))
		if len(config.ApiKeys) != 1 {
			t.Errorf("expected 1 api key, got %v", len(config.ApiKeys))
		}
		wantName, wantValue := "main-api-key", "foo"
		if config.ApiKeys[0].Name != wantName {
			t.Errorf("expected name to be %q got %q", wantName, config.ApiKeys[0].Name)
		}
		if config.ApiKeys[0].Value != wantValue {
			t.Errorf("expected value to be %q got %q", wantValue, config.ApiKeys[0].Value)
		}
	})
}
