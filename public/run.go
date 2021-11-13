package public

import (
	"fmt"
	"net/http"

	"github.com/martindrlik/play/sequence"
)

func Run(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, sequence.Get(r.Context()))
}
