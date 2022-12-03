package plugin

import (
	"fmt"
	"os/exec"
)

func build(goFile, soFile string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o="+soFile, goFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to build %q as plugin: %w", goFile, err)
	}
	return nil
}
