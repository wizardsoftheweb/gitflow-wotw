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
	ErrUnableToGitInit                     = errors.New("Unable to complete git init in the current working directory")
	ErrHeadlessRepo                        = errors.New("Unable to initialize in a bare repo")
	ErrAlreadyInitialized                  = errors.New("The repo is already initialized; try again with -f")
	ErrProductionMustDifferFromDevelopment = errors.New("The production branch must differ from the development branch")
	ErrUnableToConfigure                   = errors.New("Unable to configure the repo")
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

func CheckInitialization(context *cli.Context, repo *Repository, branch string) string {
	logrus.Trace("CheckInitialization")
	if repo.HasBranchBeenConfigured(branch) && !context.Bool("force") {
		value, _ := repo.config.Option(GIT_CONFIG_READ, "gitflow", "branch", "master")
		return value
	}
	return ""
}

func CheckIfABranchSuggestionExists(suggestions []string, branches []string) string {

	logrus.Trace("CheckIfABranchSuggestionExists")
	for _, suggestion := range suggestions {
		for _, branch := range branches {
			if suggestion == branch {
				return branch
			}
		}
	}
	return suggestions[0]
}

func ConstructMasterBranchNameSuggestions(context *cli.Context, repo Repository) error {

	logrus.Trace("ConstructMasterBranchNameSuggestions")
	if 0 == len(repo.localBranches) {
		logrus.Info("No branches exist; creating now")
		ActiveCommandInitState.MasterExistenceCheck = false
		value, err := repo.config.Option(
			GIT_CONFIG_READ,
			"gitflow",
			"branch",
			"master",
		)
		if nil != err && ErrConfigOptionNotFound != err {
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

func BuildMasterBranch(context *cli.Context, repo *Repository) string {

	logrus.Trace("BuildMasterBranch")
	master := PromptForBranchName(
		fmt.Sprintf("Branch name for prod [%s]", ActiveCommandInitState.MasterDefaultSuggestion),
	)
	repo.config.Option(GIT_CONFIG_UPDATE, "gitflow", "branch", "master", master)
	if ActiveCommandInitState.MasterExistenceCheck {
		if !repo.DoesBranchExistLocally(master) {
			remoteMaster := fmt.Sprintf("origin/%s", master)
			if repo.DoesBranchExistRemotely(remoteMaster) {
				execute("git", "branch", master, remoteMaster)
			}
		} else {
			logrus.Warning(fmt.Sprintf("The chosen master branch %s does not exist locally", master))
		}
	}
	return master
}

func ConstructDevBranchNameSuggestions(context *cli.Context, repo Repository) error {

	logrus.Trace("ConstructDevBranchNameSuggestions")
	if 0 == len(repo.localBranches) {
		logrus.Info("No branches exist; creating now")
		ActiveCommandInitState.DevExistenceCheck = false
		value, err := repo.config.Option(
			GIT_CONFIG_READ,
			"gitflow",
			"branch",
			"dev",
		)
		if nil != err && ErrConfigOptionNotFound != err {
			logrus.Fatal(err)
		}
		if "" != value {
			ActiveCommandInitState.DevDefaultSuggestion = value
		} else {
			ActiveCommandInitState.DevDefaultSuggestion = DefaultGitflowBranchDevelopmentOption.Value
		}
	} else {
		logrus.Info("What branch should be used for development releases?")
		ActiveCommandInitState.DevExistenceCheck = true
		ActiveCommandInitState.DevDefaultSuggestion = CheckIfABranchSuggestionExists(DefaultDevBranchGuesses, repo.localBranches)
	}
	return nil
}

func BuildDevBranch(context *cli.Context, repo *Repository, master string) string {

	logrus.Trace("BuildDevBranch")
	dev := PromptForBranchName(
		fmt.Sprintf("Branch name for dev [%s]", ActiveCommandInitState.DevDefaultSuggestion),
	)
	repo.config.Option(GIT_CONFIG_UPDATE, "gitflow", "branch", "dev", dev)
	logrus.Debug(repo.config.Option(GIT_CONFIG_READ, "gitflow", "branch", "dev"))
	if ActiveCommandInitState.DevExistenceCheck {
		if !repo.DoesBranchExistLocally(dev) {
			logrus.Warning(fmt.Sprintf("The chosen dev branch %s does not exist locally", dev))
		}
	}
	return dev
}

func EnsureDevAndMasterDiffer(dev string, master string) error {
	logrus.Trace("EnsureDevAndMasterDiffer")
	if dev == master {
		return ErrProductionMustDifferFromDevelopment
	}
	return nil
}

func EnsureHeadExists(context *cli.Context, master string) error {
	logrus.Trace("EnsureHeadExists")
	verify := execute("git", "rev-parse", "--quiet", "--verify", "HEAD")
	if !verify.Succeeded() {
		execute("git", "symbolic-ref", "HEAD", fmt.Sprintf("refs/heads/%s", master))
		execute("git", "commit", "--allow-empty", "--quiet", "-m", "Initial commit")
	}
	return nil
}

func EnsureDevExists(context *cli.Context, repo *Repository, dev string, master string) error {

	logrus.Trace("EnsureDevExists")
	if !repo.DoesBranchExistLocally(dev) {
		devMaster := fmt.Sprintf("origin/%s", dev)
		if repo.DoesBranchExistRemotely(devMaster) {
			execute("git", "branch", dev, devMaster)
		} else {
			execute("git", "branch", "--no-track", dev, master)
		}
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
	master := CheckInitialization(context, &repo, "master")
	if "" == master {
		if context.Bool("default") {
			master = DefaultGitflowBranchMasterOption.Value
		} else {
			ConstructMasterBranchNameSuggestions(context, repo)
			master = BuildMasterBranch(context, &repo)
		}
	}
	dev := CheckInitialization(context, &repo, "dev")
	if "" == dev {
		if context.Bool("default") {
			dev = DefaultGitflowBranchDevelopmentOption.Value
		} else {
			ConstructDevBranchNameSuggestions(context, repo)
			dev = BuildDevBranch(context, &repo, master)
		}
	}
	err = EnsureDevAndMasterDiffer(dev, master)
	if nil != err {
		log.Fatal(err)
	}
	EnsureHeadExists(context, master)
	EnsureDevExists(context, &repo, dev, master)
	repo.LoadLocalBranches()
	if !repo.HasBranchBeenConfigured("master") && !repo.HasBranchBeenConfigured("dev") {
		return ErrUnableToConfigure
	}
	execute("git", "checkout", dev)
	return nil
}
