package main

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandInit = cli.Command{
		Name: "init",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "force, f",
				Usage: "Forces a reinitialization",
			},
		},
		Action: CommandInitAction,
	}
)

func GetKeyFromRef(branch string) string {
	switch branch {
	case "master":
		return MASTER_BRANCH_KEY
	}
	return DEV_BRANCH_KEY
}

func GetValueFromRef(branch string) string {
	switch branch {
	case "master":
		return "master"
	}
	return "dev"
}

func SharedPrep(context *cli.Context, branchName string) error {
	if context.Bool("force") || !IsBranchConfigured(branchName) {
		var check_existence bool
		var suggestion string
		if 0 == len(Repo.LocalBranches()) {
			check_existence = false
			suggestion = GitConfig.GetWithDefault(GetKeyFromRef(branchName), branchName)
		} else {
			check_existence = true
			suggestion = Repo.PickGoodSuggestion(branchName)
			if "" == suggestion {
				check_existence = false
				suggestion = GitConfig.GetWithDefault(GetKeyFromRef(branchName), GetValueFromRef(branchName))
			}
		}
		newValue := PromptForInput(REF_NAME_VALIDATION, PromptMessageFromBranch(branchName, suggestion), suggestion)
		GitConfig.Write(GetKeyFromRef(branchName), newValue)
		if check_existence {
			CheckExistence(branchName, newValue)
		}
	}
	return nil
}

func CheckExistence(branchName string, newName string) {
	if "master" == branchName {
		if !Repo.HasLocalBranch(newName) {
			if Repo.HasRemoteBranch(newName) {
				ExecCmd("git", "branch", newName, fmt.Sprintf("origin/%s", newName))
			} else {
				logrus.Warn("Master", Repo.HasLocalBranch(newName))
				logrus.Warn(fmt.Sprintf("'%s'", newName))
				logrus.Error(fmt.Sprintf("Branch '%s' does not exist", newName))
				logrus.Fatal(ErrProdDoesntExist)
			}
		}
		return
	} else {
		if !Repo.HasLocalBranch(newName) {
			logrus.Fatal(ErrProdDoesntExist)
		}
		return
	}

}
func PromptMessageFromBranch(branchName string, suggestion string) string {
	switch branchName {
	case "master":
		return fmt.Sprintf(
			"Branch name for prod [%s]",
			suggestion,
		)
	default:
		return fmt.Sprintf(
			"Branch name for dev [%s]",
			suggestion,
		)
	}

}

func ParsePrefix(context *cli.Context, prefixKey string, defaultValue string) error {
	if !context.Bool("force") {
		return nil
	} else {
		var newValue string
		if context.Bool("default") {
			newValue = defaultValue
		} else {
			value := GitConfig.GetWithDefault(prefixKey, defaultValue)
			newValue = PromptForInput(
				PREFIX_NAME_VALIDATION,
				fmt.Sprintf("Prefix for %s branches? [%s]", filepath.Base(prefixKey), defaultValue),
				value,
			)
		}
		GitConfig.Write(prefixKey, newValue)
	}
	return nil
}

func InitProcedural(context *cli.Context) error {
	logrus.Trace("InitProcedural")
	if !RevParseGitDir().Succeeded() {
		GitInit()
	} else {
		if RevParseQuietVerifyHead().Succeeded() {
			IsWorkingTreeClean()
		} else {
			logrus.Fatal(ErrHeadlessRepo)
		}
	}
	logrus.Trace("Repo has be identified")
	if IsGitFlowInitialized() && !context.Bool("force") {
		logrus.Fatal(ErrAlreadyInitialized)
	}
	if context.Bool("default") {
		fmt.Fprint(context.App.Writer, "Using default branch names")
	}
	for _, branchName := range []string{"master", "dev"} {
		SharedPrep(context, branchName)
	}
	devName := GitConfig.Get(DEV_BRANCH_KEY)
	masterName := GitConfig.Get(MASTER_BRANCH_KEY)
	if devName == masterName {
		logrus.Fatal(ErrProductionMustDifferFromDevelopment)
	}
	var createdBranch bool
	if !RevParseQuietVerifyHead().Succeeded() {
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
	if !IsMasterConfigured() || !IsDevConfigured() {
		logrus.Fatal(ErrUnableToConfigure)
	}
	if createdBranch {
		ExecCmd("git", "checkout", "-q", devName)
	}

	if context.Bool("force") || !ArePrefixesConfigured() {
		fmt.Fprint(context.App.Writer, "Some prefixes need to be configured")
		for _, prefix := range DefaultPrefixes {
			ParsePrefix(context, prefix.Key, prefix.Value)
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
				logrus.Trace(prefix.Key, newValue)
				GitConfig.Write(prefix.Key, newValue)
			}
		}
	}
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Debug("CommandInitAction")
	return InitProcedural(context)
}
