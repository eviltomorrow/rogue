package pid

import (
	"fmt"
	"os"

	"github.com/eviltomorrow/rogue/lib/fs"
	"github.com/eviltomorrow/rogue/lib/runutil"
)

func CreatePidFile(path string) (func() error, error) {
	file, err := fs.CreateFlockFile(path)
	if err != nil {
		return nil, err
	}

	file.WriteString(fmt.Sprintf("%d", runutil.Pid))
	if err := file.Sync(); err != nil {
		file.Close()
		return nil, err
	}

	return func() error {
		if file != nil {
			if err := file.Close(); err != nil {
				return err
			}
			return os.Remove(path)
		}
		return fmt.Errorf("panic: pid file is nil")
	}, nil
}
