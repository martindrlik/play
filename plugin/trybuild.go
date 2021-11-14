package plugin

import (
	"log"
	"net/http"
	"os/exec"
	"path"
)

func tryBuild(rw http.ResponseWriter, goFile string) bool {
	if !tryChdir(rw, path.Dir(goFile)) {
		return false
	}
	cmd := exec.Command("go", "build", "-buildmode=plugin", goFile)
	if err := cmd.Run(); err != nil {
		log.Printf("unable to build plugin %q: %v", goFile, err)
		// note that this also might be InternalServerError
		http.Error(rw, "unable to build plugin", http.StatusBadRequest)
		return false
	}
	return true
}
