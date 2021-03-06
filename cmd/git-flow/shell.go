package gitflow

import (
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

type CommandResponse struct {
	result   string
	exitCode int
	exitErr  error
}

func (c CommandResponse) Succeeded() bool {
	return 0 == c.exitCode
}

func (c CommandResponse) Bool() bool {
	return c.Succeeded()
}

func (c CommandResponse) String() string {
	return string(c.result)
}

func parseExitCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func ExecCmd(args ...string) CommandResponse {
	logrus.Trace(args)
	process := exec.Command(args[0], args[1:]...)
	combined, err := process.CombinedOutput()
	return CommandResponse{
		result:   string(combined),
		exitCode: parseExitCode(err),
		exitErr:  err,
	}
}
