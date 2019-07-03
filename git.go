package main

import "errors"

var (
	ErrUnstagedChanges  = errors.New("There are unstaged changes in your working directory")
	ErrIndexUncommitted = errors.New("There are uncommitted changes in your index")
)

func IsRepoHeadless() bool {
	revparse := execute("git", "rev-parse", "--quiet", "--verify", "HEAD")
	return !revparse.Succeeded()
}

func IsWorkingTreeClean() bool {
	cleanliness := execute("git", "diff", "--no-ext-diff", "--ignore-submodules", "--quiet", "--exit-code")
	return cleanliness.Succeeded()
}

func IsIndexClean() bool {
	cleanliness := execute("git", "diff-index", "--cached", "--quiet", "--ignore-submodules", "HEAD")
	return cleanliness.Succeeded()
}

func EnsureCleanWorkingTree() error {
	if !IsWorkingTreeClean() {
		return ErrUnstagedChanges
	}
	if !IsIndexClean() {
		return ErrIndexUncommitted
	}
	return nil
}
