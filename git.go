package main

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var (
	ErrUnstagedChanges  = errors.New("There are unstaged changes in your working directory")
	ErrIndexUncommitted = errors.New("There are uncommitted changes in your index")
)

func IsRepoHeadless() bool {
	logrus.Trace("IsRepoHeadless")
	revparse := execute("git", "rev-parse", "--quiet", "--verify", "HEAD")
	return !revparse.Succeeded()
}

func IsWorkingTreeClean() bool {
	logrus.Trace("IsWorkingTreeClean")
	cleanliness := execute("git", "diff", "--no-ext-diff", "--ignore-submodules", "--quiet", "--exit-code")
	return cleanliness.Succeeded()
}

func IsIndexClean() bool {
	logrus.Trace("IsIndexClean")
	cleanliness := execute("git", "diff-index", "--cached", "--quiet", "--ignore-submodules", "HEAD")
	return cleanliness.Succeeded()
}

func EnsureCleanWorkingTree() error {
	logrus.Trace("EnsureCleanWorkingTree")
	if !IsWorkingTreeClean() {
		return ErrUnstagedChanges
	}
	if !IsIndexClean() {
		return ErrIndexUncommitted
	}
	return nil
}
