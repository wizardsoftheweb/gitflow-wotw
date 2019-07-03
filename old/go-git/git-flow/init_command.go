package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	format "gopkg.in/src-d/go-git.v4/plumbing/format/config"
)

const (
	MessageNoBranches = "There are no branches; everything must be created"
)

var (
	BranchNameSuggestions = map[string][]string{
		"master": []string{
			"master",
			"prod",
			"production",
			"main",
		},
		"dev": []string{
			"dev",
			"development",
		},
	}
	GitflowBranchOptArgs = map[string]ConfigOptionArgs{
		"master": GitflowBranchMasterOption,
		"dev":    GitflowBranchDevelopmentOption,
	}
)

type CommandInitOptions struct {
	HasBeenInitialized           bool
	ShouldCheckExistence         bool
	DefaultMasterSuggestion      string
	DefaultDevelopmentSuggestion string
	LocalBranches                []string
	FinalMasterValue             string
	FinalDevelopmentValue        string
}

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

func (init_options *CommandInitOptions) GenerateBranchSuggestion(gitflow_branch string) string {
	suggestions, ok := BranchNameSuggestions[gitflow_branch]
	if !ok {
		logrus.Warn("%s does not have any suggestions", gitflow_branch)
		return ""
	}
	for _, name := range suggestions {
		if StringSliceContains(init_options.LocalBranches, name) {
			return name
		}
	}
	return ""
}

func (init_options *CommandInitOptions) ConfigureMasterBranch(context *cli.Context, repo_config *config.Config) error {
	if init_options.HasBeenInitialized && !context.Bool("force") {
		init_options.FinalMasterValue = GitflowBranchMasterOption.getValueWithDefault(repo_config.Raw, false)
	} else {
		if 1 > len(init_options.LocalBranches) {
			fmt.Fprintln(context.App.Writer, MessageNoBranches)
			init_options.ShouldCheckExistence = false
			init_options.DefaultMasterSuggestion = GitflowBranchMasterOption.getValueWithDefault(repo_config.Raw, true)
		} else {
			init_options.ShouldCheckExistence = true
			chosen_name := PromptForBranchName(fmt.Sprintf(
				"Branch for production releases [%s]",
				init_options.GenerateBranchSuggestion("master"),
			))
			if "" == chosen_name {
				chosen_name = GitflowBranchMasterOption.getValueWithDefault(repo_config.Raw, true)
			}
			init_options.FinalMasterValue = chosen_name
		}
	}
	return nil
}

func GitFlowInit(context *cli.Context) error {
	repo_path, _ := os.Getwd()
	logrus.Debug("GitFlowInit")
	repo, err := EnsureRepoIsUsable(repo_path)
	CheckError(err)
	repo_config, err := LoadConfig(repo)
	CheckError(err)
	init_options := &CommandInitOptions{
		HasBeenInitialized: EnsureNecessaryInitOptionsAreSet(repo_config.Raw),
		LocalBranches:      GetLocalBranchNames(repo_config),
	}
	err = init_options.ConfigureMasterBranch(context, repo_config)
	if nil != err {
		logrus.Warn(err)
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
