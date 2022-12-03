package plugin

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/martindrlik/play/metrics"
)

var (
	MaxUploadFileLength = 16e3
)

// Upload creates handler that produces uploaded content by using produce func.
func Upload(produce func(value, key []byte) error) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.ContentLength > int64(MaxUploadFileLength) {
			metrics.UploadMaxFileLengthExceeded()
			http.Error(
				rw,
				fmt.Sprintf("max upload file length (%v B) exceeded", MaxUploadFileLength),
				http.StatusRequestEntityTooLarge)
			return
		}
		name, ok := tryGetName(rw, r)
		if !ok {
			return
		}
		value := &bytes.Buffer{}
		_, err := io.CopyN(value, r.Body, int64(MaxUploadFileLength))
		if err != nil && err != io.EOF {
			metrics.UploadReadingBodyError()
			log.Printf("upload reading body: %v", err)
			http.Error(rw, "unable to read request body", http.StatusInternalServerError)
			return
		}
		err = produce(value.Bytes(), []byte(name))
		if err != nil {
			metrics.UploadStoringError()
			log.Printf("upload storing body: %v", err)
			http.Error(rw, "unable to store uploaded content", http.StatusInternalServerError)
			return
		}
	}
}

func tryGetName(rw http.ResponseWriter, r *http.Request) (name string, ok bool) {
	name = r.URL.Path[len("/upload"):]
	if name == "" {
		http.Error(rw, "usage /upload/example", http.StatusBadRequest)
		return "", false
	}
	return name, true
}
