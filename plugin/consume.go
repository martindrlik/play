package plugin

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

// Consume compiles handler's source code given by content and makes it accessible as an API given by path.
func Consume(content []byte, path string) {
	main, err := consume(content, path)
	storageMutex.Lock()
	defer storageMutex.Unlock()
	if err == nil {
		// set new version only if there is no error
		plugins[path] = main
	}
	analyze[path] = err
}

func consume(content []byte, path string) (func(http.ResponseWriter, *http.Request), error) {
	dir, goFile, soFile, err := getFileNames(path)
	if err != nil {
		return nil, fmt.Errorf("unable to get names: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create dir: %w", err)
	}
	if err := writeGoFile(content, goFile); err != nil {
		return nil, fmt.Errorf("unable to write %q: %w", goFile, err)
	}
	if err := build(goFile); err != nil {
		return nil, fmt.Errorf("unable to build %q: %w", goFile, err)
	}
	main, err := tryLookupHandler(soFile)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup %q: %w", soFile, err)
	}
	return main, nil
}

func getFileNames(name string) (dir, goFile, soFile string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("unable to get working directory: %w", err)
		return
	}
	s := path.Join(wd, name)
	dir = path.Dir(s)
	goFile = s + ".go"
	soFile = s + ".so"
	return
}

func writeGoFile(content []byte, goFile string) error {
	f, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("unable to create file %q: %w", goFile, err)
	}
	defer f.Close()
	n, err := f.Write(content)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", goFile, err)
	}
	if n != len(content) {
		return fmt.Errorf("unexpected number of bytes %v written (expected %v)", n, len(content))
	}
	log.Printf("file %q (%v B) created 🎉", goFile, n)
	return nil
}
