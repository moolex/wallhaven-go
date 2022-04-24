//go:build darwin

package wpset

import (
	"fmt"
	"os/exec"
)

func Image(path string) error {
	return exec.Command(
		"/usr/bin/osascript",
		"-e",
		fmt.Sprintf(`tell application "Finder" to set desktop picture to POSIX file "%s"`, path),
	).Run()
}
