package plugin

import (
	"fmt"
	"os/exec"
)

func build(goFile string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", goFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to build %q as plugin: %w", goFile, err)
	}
	return nil
}
