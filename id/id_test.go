package id_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martindrlik/play/id"
)

func TestId(t *testing.T) {
	rid := ""
	id.Gen()(func(rw http.ResponseWriter, r *http.Request) {
		rid = id.Get(rw)
	})(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil))
	if rid == "" {
		t.Error("expected rid to have value got empty")
	}
}
