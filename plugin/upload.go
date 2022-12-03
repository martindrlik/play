package plugin

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	MaxUploadFileLength = 16e3
)

// Upload creates handler that produces uploaded content by using produce func.
func Upload(produce func(value, key []byte) error) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.ContentLength > int64(MaxUploadFileLength) {
			http.Error(rw, "max upload file length exceeded", http.StatusBadRequest)
			return
		}
		name, ok := tryGetName(rw, r)
		if !ok {
			return
		}
		value := &bytes.Buffer{}
		_, err := io.CopyN(value, r.Body, int64(MaxUploadFileLength))
		if err != nil {
			http.Error(rw, "unable to read request body", http.StatusInternalServerError)
			return
		}
		err = produce(value.Bytes(), []byte(name))
		if err != nil {
			http.Error(rw, "unable to store request body", http.StatusInternalServerError)
			return
		}
	}
}

// Upload allows uploading go file, builds it as a plugin and returns url path that runs it.
func x(rw http.ResponseWriter, r *http.Request) {
	if r.ContentLength > int64(MaxUploadFileLength) {
		http.Error(rw, "max upload file length exceeded", http.StatusBadRequest)
		return
	}

	name, ok := tryGetName(rw, r)
	if !ok {
		return
	}
	dir, goFile, soFile, ok := tryGetDir(rw, name)
	if !ok {
		return
	}
	if !tryMakeDirAll(rw, dir) {
		return
	}
	if !tryUpload(rw, r, goFile) {
		return
	}
	if !tryBuild(rw, goFile) {
		return
	}
	main, ok := tryLookup(rw, soFile)
	if !ok {
		return
	}

	pluginsMutex.Lock()
	defer pluginsMutex.Unlock()
	plugins[name] = main
}

func tryGetName(rw http.ResponseWriter, r *http.Request) (name string, ok bool) {
	name = r.URL.Path[len("/upload"):]
	if name == "" {
		http.Error(rw, "usage /upload/example", http.StatusBadRequest)
		return "", false
	}
	return name, true
}

func tryGetDir(rw http.ResponseWriter, name string) (dir, goFile, soFile string, ok bool) {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("unable to get working directory: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return "", "", "", false
	}
	s := path.Join(wd, name)
	dir = path.Dir(s)
	goFile = s + ".go"
	soFile = s + ".so"
	ok = true
	return
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
