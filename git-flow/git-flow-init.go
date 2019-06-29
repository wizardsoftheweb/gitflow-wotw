package main

import (
	"errors"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func CheckError(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

func IsCwdInRepo() bool {
	current_path, _ := os.Getwd()
	_, err := OpenRepoFromPath(current_path)
	if git.ErrRepositoryNotExists == err {
		return false
	}
	return true
}

func CanAppInitRepo(repo_path string) bool {
	err := GitInit(repo_path)
	if git.ErrRepositoryAlreadyExists == err {
		return false
	} else {
		CheckError(err)
	}
	return true
}

func EnsureRepoIsUsable(repo_path string) error {
	if !CanAppInitRepo(repo_path) {
		repo = OpenRepoFromPath(repo_path)
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
