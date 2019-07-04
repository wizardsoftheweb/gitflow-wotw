package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	ScreenSizePattern = regexp.MustCompile(`^.*?(\d+)\.*?(\d+)\.*?$`)
)

var (
	SizeOfScreen ScreenSize
)

type CommandResponse struct {
	result   string
	exitCode int
	exitErr  error
}

func (response CommandResponse) Succeeded() bool {
	return 0 == response.exitCode
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

type ScreenSize struct {
	Width  int
	Height int
}

// https://stackoverflow.com/a/16569795
func GetTermSize() ScreenSize {
	logrus.Debug("GetTermSize")
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	screenSize := [2]int{}
	for index, match := range ScreenSizePattern.FindAllStringSubmatch(string(out), -1) {
		screenSize[index], _ = strconv.Atoi(strings.Join(match, ""))
	}
	return ScreenSize{Width: screenSize[0], Height: screenSize[1]}
}
