package plugin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/martindrlik/play/sequence"
)

var (
	MaxUploadFileLength = 16e3
)

// Upload allows uploading go file, builds it as a plugin and returns url path that runs it.
func Upload(rw http.ResponseWriter, r *http.Request) {
	if r.ContentLength > int64(MaxUploadFileLength) {
		http.Error(rw, "max upload file length exceeded", http.StatusBadRequest)
		return
	}
	name := uploadName(sequence.Get(r.Context()))
	dir := path.Dir(name)
	if !tryMakeDirAll(rw, dir) {
		return
	}
	goFile := name + ".go"
	if !tryUpload(rw, r, goFile) {
		return
	}
	if !tryBuild(rw, goFile) {
		return
	}
	soFile := name + ".so"
	main, ok := tryLookup(rw, soFile)
	if !ok {
		return
	}

	key := path.Base(name)
	func() {
		pluginsMutex.Lock()
		defer pluginsMutex.Unlock()
		plugins[key] = main
	}()
	fmt.Fprintf(rw, "/run/%s", key)
}

func uploadName(seq int64) string {
	return path.Join(
		os.Getenv("GOPATH"),
		"src/github.com/martindrlik/play/plugins", // TODO better path
		strconv.FormatInt(seq, 36))
}

func tryMakeDirAll(rw http.ResponseWriter, dir string) bool {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("unable to create directory %q: %v", dir, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

func tryUpload(rw http.ResponseWriter, r *http.Request, goFile string) bool {
	f, err := os.Create(goFile)
	if err != nil {
		log.Printf("unable to create file %q: %v", goFile, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return false
	}
	written, err := io.CopyN(f, r.Body, int64(MaxUploadFileLength))
	if err != nil && err != io.EOF {
		log.Printf("unable to copy request body: %v", err)
		// TODO BadRequest or InternalServerError
		http.Error(rw, "unable to upload file", http.StatusInternalServerError)
		return false
	}
	log.Printf("file %q (%d) uploaded", goFile, written)
	return true
}

func tryChdir(rw http.ResponseWriter, dir string) bool {
	// TODO there should be no need for changing working directory
	cur, err := os.Getwd()
	if err != nil {
		log.Printf("unable to get working directory: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if cur != dir {
		err = os.Chdir(dir)
	}
	if err != nil {
		log.Printf("unable to change working directory to %q: %v", dir, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
