package main

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

func (response *CommandResponse) Succeeded() bool {
	return 0 == response.exitCode
}

func parseExitCode(err error) int {
	logrus.Trace("parseExitCode")
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func ExecCmd(args ...string) CommandResponse {
	logrus.Trace("ExecCmd")
	logrus.Debug(args)
	process := exec.Command(args[0], args[1:]...)
	combined, err := process.CombinedOutput()
	CheckError(err)
	return CommandResponse{
		result:   string(combined),
		exitCode: parseExitCode(err),
		exitErr:  err,
	}
}
