package fs

import (
	"fmt"
	"os"
)

func CreateDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return os.MkdirAll(dir, 0755)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a dir, path: %v", dir)
	}
	return nil
}
