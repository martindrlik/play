package measure_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martindrlik/play/measure"
)

func TestMeasure(t *testing.T) {
	called := 0
	observe := func(f float64) { called++ }
	measure.Measure(observe, func(w http.ResponseWriter, r *http.Request) {})(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodPost, "/", nil))
	if called != 1 {
		t.Errorf("expected observe to be called once, called %v", called)
	}
}
