package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

type CommandResponse struct {
	stdout   string
	stderr   string
	exitCode int
	exitErr  error
}

func (response *CommandResponse) Succeeded() bool {
	return 0 == response.exitCode
}

// https://stackoverflow.com/a/10385867
func parseExitCode(err error) int {
	logrus.Trace("parseExitCode")
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func execute(command ...string) CommandResponse {
	logrus.Trace("execute")
	logrus.Debug("execute", command)
	stdout, stderr, err := RunCommand(command)
	return CommandResponse{
		stdout:   stdout,
		stderr:   stderr,
		exitCode: parseExitCode(err),
		exitErr:  err,
	}
}

func RunCommand(sanitized_command []string) (string, string, error) {
	logrus.Trace("RunCommand")
	command := exec.Command(sanitized_command[0], sanitized_command[1:]...)
	stdout, err := command.StdoutPipe()
	if nil != err {
		log.Fatal(err)
	}
	stderr, err := command.StderrPipe()
	if nil != err {
		log.Fatal(err)
	}
	err = command.Start()
	if nil != err {
		log.Fatal(err)
	}
	final_out, _ := ioutil.ReadAll(stdout)
	final_err, _ := ioutil.ReadAll(stderr)
	err = command.Wait()
	logrus.Debug(string(final_out))
	logrus.Debug(string(final_err))
	return string(final_out), string(final_err), err
}
