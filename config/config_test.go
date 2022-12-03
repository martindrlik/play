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
}
