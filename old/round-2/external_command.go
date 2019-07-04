package main

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Execute struct {
	Command []string
}

type Git struct {
	Execute
}

type RevParse struct {
	Execute
}

func (e *Execute) exec() error {
	process := exec.Command(e.Command[0], e.Command[1:]...)
	combined, err := process.CombinedOutput()
	if nil != err {
		logrus.Fatal(err)
	}
	fmt.Println(string(combined))
	return nil
}

func (g *Git) RevParse() *RevParse {
	revparse := &RevParse{}
	revparse.Command = append(g.Command, "rev-parse")
	return revparse
}

func (r *RevParse) GitDir() *RevParse {
	r.Command = append(r.Command, "--git-dir")
	return r
}

func TestRun() {
	git := &Git{}
	git.Command = []string{"git"}
	git.RevParse().GitDir().exec()
}
