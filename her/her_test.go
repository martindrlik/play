package her_test

import (
	"errors"
	"testing"

	"github.com/martindrlik/play/her"
)

func TestMust(t *testing.T) {
	defer func() {
		x := recover()
		if err, ok := x.(error); !ok || err.Error() != "foo" {
			t.Errorf("expected panic error(foo) got %T(%v)", x, x)
		}
	}()
	if x := her.Must(1, nil); x != 1 {
		t.Errorf("expected 1 got %v", x)
	}
	her.Must(1, errors.New("foo"))
}

func TestMust1(t *testing.T) {
	her.Must1(nil) // expecting no panic
	defer func() {
		x := recover()
		if err, ok := x.(error); !ok || err.Error() != "foo" {
			t.Errorf("expected panic error(foo) got %T(%v)", x, x)
		}
	}()
	her.Must1(errors.New("foo"))
}
