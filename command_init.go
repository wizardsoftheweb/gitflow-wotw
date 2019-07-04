package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandInit = cli.Command{
		Name:   "release",
		Flags:  []cli.Flag{},
		Action: CommandInitAction,
	}
)

func InitProcedural(context *cli.Context) error {
	result := RevParseGitDir()
	if !result.Succeeded() {
		GitInit()
	} else {
		result = RevParseQuietVerifyHead()
		if !result.Succeeded() {
			IsWorkingTreeClean()
		} else {
			logrus.Fatal(ErrHeadlessRepo)
		}
	}
	if IsGitFlowInitialized() && !context.Bool("force") {
		logrus.Fatal(ErrAlreadyInitialized)
	}
	if context.Bool("default") {
		fmt.Fprint(context.App.Writer, "Using default branch names\n")
	}
	var masterName string
	var check_existence bool
	var suggestion string
	if IsMasterConfigured() && !context.Bool("force") {
		masterName = GitConfig.Get(MASTER_BRANCH_KEY)
	} else {
		if 0 == len(Repo.LocalBranches()) {
			fmt.Fprint(context.App.Writer, "No branches exist; creating now\n")
			check_existence = false
			value := GitConfig.Get(MASTER_BRANCH_KEY)
			if "" == value {
				suggestion = "master"
			} else {
				suggestion = value
			}
		} else {
			fmt.Fprint(context.App.Writer, "Which branch should be used for production?\n")
			localBranches := Repo.LocalBranches()
			for _, branch := range localBranches {
				fmt.Fprintf(context.App.Writer, "\t-%s", branch)
			}
			check_existence = true
			suggestion = Repo.PickGoodMasterSuggestion()
		}
		masterName := PromptForBranchName(
			fmt.Sprintf(
				"Branch name for prod [%s]",
				suggestion,
			),
			suggestion,
		)
		if check_existence {
			if !Repo.HasLocalBranch(masterName) {
				if Repo.HasRemoteBranch(masterName) {
					ExecCmd("git", "branch", masterName, fmt.Sprintf("origin/%s", masterName))
				}
			} else {
				logrus.Error(fmt.Sprintf("Branch '%s' does not exist"))
				logrus.Fatal(ErrProdDoesntExist)
			}
		}
	}
	logrus.Debug(masterName)
	logrus.Debug(check_existence)
	logrus.Debug(suggestion)
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Debug("CommandInitAction")
	return nil
}
