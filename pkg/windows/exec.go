package windows

import (
	"bytes"
	"os/exec"
)

// RunCommand выполняет системную команду и возвращает вывод
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}
