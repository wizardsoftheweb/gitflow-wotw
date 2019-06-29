package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	format "gopkg.in/src-d/go-git.v4/plumbing/format/config"
)

func IsCwdInRepo(repo_path string) bool {
	logrus.Debug("IsCwdInRepo")
	_, err := OpenRepoFromPath(repo_path)
	if git.ErrRepositoryNotExists == err {
		return false
	}
	return true
}

func CanAppInitRepo(repo_path string) bool {
	logrus.Debug("CanAppInitRepo")
	if IsCwdInRepo(repo_path) {
		return false
	}
	err := GitInit(repo_path)
	CheckError(err)
	return true
}

func EnsureRepoIsUsable(repo_path string) (*git.Repository, error) {
	logrus.Debug("EnsureRepoIsUsable")
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

var (
	NecessaryInitSettings = []ConfigOptionArgs{
		GitflowBranchDevelopmentOption,
		GitflowBranchMasterOption,
		GitflowPrefixFeatureOption,
		GitflowPrefixHotfixOption,
		GitflowPrefixReleaseOption,
		GitflowPrefixSupportOption,
		GitflowPrefixVersiontagOption,
	}
)

func EnsureNecessaryInitOptionsAreSet(git_config *format.Config) bool {
	logrus.Debug("EnsureNecessaryInitOptionsAreSet")
	for _, option_arg_set := range NecessaryInitSettings {
		logrus.Trace(option_arg_set)
		if !option_arg_set.isOptionSetInConfig(git_config) {
			return false
		}
	}
	return true

}

func EnsureBranchesExist() {
	println("rad")
}

func GitFlowInit(context *cli.Context) error {
	repo_path, _ := os.Getwd()
	logrus.Debug("GitFlowInit")
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

var (
	CliFlagForce = cli.BoolFlag{
		Name: "force",
	}
	CommandInit = cli.Command{
		Name: "init",
		Flags: []cli.Flag{
			CliFlagForce,
		},
		Action: GitFlowInit,
	}
)
