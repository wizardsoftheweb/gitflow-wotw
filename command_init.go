package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type CommandInitState struct {
	DevDefaultSuggestion    string
	DevExistenceCheck       bool
	MasterDefaultSuggestion string
	MasterExistenceCheck    bool
}

var ActiveCommandInitState CommandInitState

var (
	ErrUnableToGitInit    = errors.New("Unable to complete git init in the current working directory")
	ErrHeadlessRepo       = errors.New("Unable to initialize in a bare repo")
	ErrAlreadyInitialized = errors.New("The repo is already initialized; try again with -f")
)

var (
	DefaultMasterBranchGuesses = []string{"master", "prod", "production", "main"}
	DefaultDevBranchGuesses    = []string{"dev", "development"}
)

var (
	CommandInit = cli.Command{
		Name: "init",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "force, f",
				Usage: "Forces the setting to be reinitialized",
			},
			cli.BoolFlag{
				Name:  "defaults, d",
				Usage: "Applies the defaults without prompting (when available)",
			},
		},
		Action: CommandInitAction,
	}
)

func EnsureRepoIsAvailable(directory string) (Repository, error) {
	logrus.Trace("EnsureRepoIsAvailable")
	var repo Repository
	var err error
	result := execute("git", "rev-parse", "--git-dir")
	if !result.Succeeded() {
		result = execute("git", "init")
		if !result.Succeeded() {
			return repo, ErrUnableToGitInit
		}
	} else {
		headless := IsRepoHeadless()
		if !headless {
			err = EnsureCleanWorkingTree()
			if nil != err {
				return repo, err
			}
		}

	}
	repo.LoadOrInit(directory)
	return repo, nil
}

func CheckInitialization(context *cli.Context, repo Repository) error {
	logrus.Trace("CheckInitialization")
	if IsGitflowInitialized(repo) && context.Bool("force") {
		return ErrAlreadyInitialized
	}
	return nil
}

func CheckIfABranchSuggestionExists(suggestions []string, branches []string) string {
	for _, suggestion := range suggestions {
		for _, branch := range branches {
			if suggestion == branch {
				return branch
			}
		}
	}
	return suggestions[0]
}

func ConstructBranchNameSuggestions(context *cli.Context, repo Repository) error {
	if 0 == len(repo.localBranches) {
		logrus.Info("No branches exist; creating now")
		ActiveCommandInitState.MasterExistenceCheck = false
		value, err := repo.config.Option(
			GIT_CONFIG_READ,
			"gitflow",
			"branch",
			"master",
		)
		if nil != err {
			logrus.Fatal(err)
		}
		if "" != value {
			ActiveCommandInitState.MasterDefaultSuggestion = value
		} else {
			ActiveCommandInitState.MasterDefaultSuggestion = DefaultGitflowBranchMasterOption.Value
		}
	} else {
		logrus.Info("What branch should be used for production?")
		ActiveCommandInitState.MasterExistenceCheck = true
		ActiveCommandInitState.MasterDefaultSuggestion = CheckIfABranchSuggestionExists(DefaultMasterBranchGuesses, repo.localBranches)
	}
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Trace("CommandInitAction")
	ActiveCommandInitState = CommandInitState{
		DevDefaultSuggestion:    DefaultGitflowBranchDevelopmentOption.Value,
		DevExistenceCheck:       false,
		MasterDefaultSuggestion: DefaultGitflowBranchMasterOption.Value,
		MasterExistenceCheck:    false,
	}
	directory, _ := os.Getwd()
	repo, err := EnsureRepoIsAvailable(directory)
	if nil != err {
		log.Fatal(err)
	}
	if context.Bool("default") {
		logrus.Info("Using default branches")
	}
	repo.LoadLocalBranches()
	fmt.Println(repo)
	execute("pwd")
	return nil
}
