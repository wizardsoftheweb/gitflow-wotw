package main

import (
	"errors"

	"gopkg.in/src-d/go-git.v4"
)

func IsCwdInRepo(repo_path string) bool {
	_, err := OpenRepoFromPath(repo_path)
	if git.ErrRepositoryNotExists == err {
		return false
	}
	return true
}

func CanAppInitRepo(repo_path string) bool {
	if IsCwdInRepo(repo_path) {
		return false
	}
	err := GitInit(repo_path)
	CheckError(err)
	return true
}

func EnsureRepoIsUsable(repo_path string) error {
	if !CanAppInitRepo(repo_path) {
		repo, err := OpenRepoFromPath(repo_path)
		CheckError(err)
		if IsRepoHeadless(repo) {
			return errors.New("This repo does not have a proper HEAD")
		} else if AreThereUnstagedChanges(repo, true) {
			return errors.New("There are unstaged changes")
		}
	}
	return nil
}

func GitFlowInit(repo_path string) {
	err := EnsureRepoIsUsable(repo_path)
	CheckError(err)
}
