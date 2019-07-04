package main

import "github.com/sirupsen/logrus"

func IsWorkingTreeClean() bool {
	logrus.Debug("IsWorkingTreeClean")
	result := ExecCmd("git", "diff", "--no-ext-diff", "--ignore-submodules", "--quiet", "--exit-code")
	if !result.Succeeded() {
		logrus.Fatal(ErrUnstagedChanges)
	}
	result = ExecCmd("git", "diff-index", "--cached", "--quiet", "--ignore-submodules", "HEAD", "--")
	if !result.Succeeded() {
		logrus.Fatal(ErrIndexUncommitted)
	}
	return true
}
