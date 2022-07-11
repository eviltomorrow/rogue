package fs

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func CreateFlockFile(path string) (*os.File, error) {
	return createFlockFile(path)
}

func createFlockFile(flockFile string) (*os.File, error) {
	flockF, err := os.Create(flockFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create lock file %q: %w", flockFile, err)
	}
	if err := unix.Flock(int(flockF.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		return nil, fmt.Errorf("cannot acquire lock on file %q: %w", flockFile, err)
	}
	return flockF, nil
}
