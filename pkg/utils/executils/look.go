package executils

import (
	"fmt"
	"os/exec"
)

func HostHasBinaries(binaries ...string) error {
	for _, binary := range binaries {
		_, err := exec.LookPath(binary)
		if err != nil {
			return fmt.Errorf("look path for '%s': %w", binary, err)
		}
	}

	return nil
}
