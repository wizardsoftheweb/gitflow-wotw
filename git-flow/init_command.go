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

func EnsureRepoIsUsable(repo_path string) (*git.Repository, error) {
	CanAppInitRepo(repo_path)
	repo, err := OpenRepoFromPath(repo_path)
	CheckError(err)
	if IsRepoHeadless(repo) {
		return nil, errors.New("This repo does not have a proper HEAD")
	} else if AreThereUnstagedChanges(repo, true) {
		return nil, errors.New("There are unstaged changes")
	}
	return repo, nil
}

func GitFlowInit(repo_path string) error {
	repo, err := EnsureRepoIsUsable(repo_path)
	CheckError(err)
	config, err := LoadConfig(repo)
	CheckError(err)
	init_options := EnsureNecessaryInitOptionsAreSet(config.Raw)
	if !init_options {
		return errors.New("Whoops")
	}
	return nil
}
