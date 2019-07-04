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
	if IsMasterConfigured() && !context.Bool("force") {
		masterName = GitConfig.Get(MASTER_BRANCH_KEY)
	} else {
		var check_existence bool
		var suggestion string
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
		masterName = PromptForInput(
			REF_NAME_VALIDATION,
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
		GitConfig.Write(MASTER_BRANCH_KEY, masterName)
	}
	var devName string
	if IsDevConfigured() && !context.Bool("force") {
		devName = GitConfig.Get(DEV_BRANCH_KEY)
	} else {

		var check_existence bool
		var suggestion string
		if 0 == len(Repo.LocalBranches()) {
			check_existence = false
			value := GitConfig.Get(DEV_BRANCH_KEY)
			if "" == value {
				suggestion = "dev"
			} else {
				suggestion = value
			}
		} else {
			fmt.Fprint(context.App.Writer, "Which branch should be used for development?\n")
			localBranches := Repo.LocalBranches()
			for _, branch := range localBranches {
				if masterName == branch {
					continue
				}
				fmt.Fprintf(context.App.Writer, "\t-%s", branch)
			}
			check_existence = true
			suggestion = Repo.PickGoodDevSuggestion(masterName)
			if "" == suggestion {
				check_existence = false
				value := GitConfig.Get(DEV_BRANCH_KEY)
				if "" == value {
					suggestion = "dev"
				} else {
					suggestion = value
				}
			}
		}
		devName := PromptForInput(
			REF_NAME_VALIDATION,
			fmt.Sprintf(
				"Branch name for dev [%s]",
				suggestion,
			),
			suggestion,
		)
		if devName == masterName {
			logrus.Fatal(ErrProductionMustDifferFromDevelopment)
		}
		if check_existence {
			if !Repo.HasLocalBranch(devName) {
				logrus.Fatal(ErrProdDoesntExist)
			}
		}
	}
	result = RevParseQuietVerifyHead()
	var createdBranch bool
	if !result.Succeeded() {
		ExecCmd("git", "symbolic-ref", "HEAD", fmt.Sprintf("refs/heads/%s", masterName))
		ExecCmd("git", "commit", "--allow-empty", "--quiet", "-m", "Initial commit")
		createdBranch = true
	}
	if !Repo.HasLocalBranch(devName) {
		if Repo.HasRemoteBranch(devName) {
			ExecCmd("git", "branch", devName, fmt.Sprintf("origin/%s", devName))
		} else {
			ExecCmd("git", "branch", "--no-track", devName, masterName)
		}
		createdBranch = true
	}
	if !IsGitFlowInitialized() {
		logrus.Fatal(ErrUnableToConfigure)
	}
	if createdBranch {
		ExecCmd("git", "checkout", "-q", devName)
	}

	if context.Bool("force") || !ArePrefixesConfigured() {
		fmt.Fprint(context.App.Writer, "Some prefixes need to be configured\n")
		for _, prefix := range DefaultPrefixes {
			value := GitConfig.Get(prefix.Key)
			defaultValue := value
			if context.Bool("force") || "" == value {
				var newValue string
				if context.Bool("default") {
					newValue = prefix.Value
				} else {
					if "" == value {
						defaultValue = prefix.Value
					}
					newValue = PromptForInput(
						PREFIX_NAME_VALIDATION,
						fmt.Sprintf("Prefix for %s branches? [%s]", prefix.Value, defaultValue),
						defaultValue,
					)
				}
				GitConfig.Write(prefix.Key, newValue)
			}
		}
		for _, prefix := range DefaultTags {
			value := GitConfig.Get(prefix.Key)
			defaultValue := value
			if context.Bool("force") || "" == value {
				var newValue string
				if context.Bool("default") {
					newValue = prefix.Value
				} else {
					if "" == value {
						defaultValue = prefix.Value
					}
					newValue = PromptForInput(
						TAG_NAME_VALIDATION,
						fmt.Sprintf("Prefix for %s tags? [%s]", prefix.Value, defaultValue),
						defaultValue,
					)
				}
				GitConfig.Write(prefix.Key, newValue)
			}
		}
	}
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Debug("CommandInitAction")
	return nil
}
