package procutil

import (
	"io"
	"os"
	"os/exec"
)

var (
	HomeDir = "."
)

func RunInBackground(name string, args []string, reader io.Reader, writer io.WriteCloser) error {
	var data = make([]string, 0, len(args)+1)
	data = append(data, name)
	data = append(data, args...)
	var cmd = &exec.Cmd{
		Path:   "/proc/self/exe",
		Args:   data,
		Stdout: writer,
		Stderr: writer,
		Stdin:  reader,
	}
	cmd.Dir = HomeDir
	cmd.Env = os.Environ()

	return cmd.Start()
}
