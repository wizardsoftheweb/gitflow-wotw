package gitflow

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
		if 0 == len(main.LocalBranches()) {
			check_existence = false
			suggestion = GetWithDefault(GetKeyFromRef(branchName), branchName)
		} else {
			check_existence = true
			suggestion = main.PickGoodSuggestion(branchName)
			if "" == suggestion {
				check_existence = false
				suggestion = GetWithDefault(GetKeyFromRef(branchName), GetValueFromRef(branchName))
			}
		}
		newValue := PromptForInput(REF_NAME_VALIDATION, PromptMessageFromBranch(branchName, suggestion), suggestion)
		Write(GetKeyFromRef(branchName), newValue)
		if check_existence {
			CheckExistence(branchName, newValue)
		}
	}
	return nil
}

func CheckExistence(branchName string, newName string) {
	if "master" == branchName {
		if !main.HasLocalBranch(newName) {
			if main.HasRemoteBranch(newName) {
				ExecCmd("git", "branch", newName, fmt.Sprintf("origin/%s", newName))
			} else {
				logrus.Warn("Master", main.HasLocalBranch(newName))
				logrus.Warn(fmt.Sprintf("'%s'", newName))
				logrus.Error(fmt.Sprintf("Branch '%s' does not exist", newName))
				logrus.Fatal(ErrProdDoesntExist)
			}
		}
		return
	} else {
		if !main.HasLocalBranch(newName) {
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
			value := GetWithDefault(prefixKey, defaultValue)
			newValue = PromptForInput(
				PREFIX_NAME_VALIDATION,
				fmt.Sprintf("Prefix for %s branches? [%s]", filepath.Base(prefixKey), defaultValue),
				value,
			)
		}
		Write(prefixKey, newValue)
	}
	return nil
}

func InitProcedural(context *cli.Context) error {
	if !main.Succeeded() {
		GitInit()
	} else {
		if main.Succeeded() {
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
	devName := Get(DEV_BRANCH_KEY)
	masterName := Get(MASTER_BRANCH_KEY)
	if devName == masterName {
		logrus.Fatal(ErrProductionMustDifferFromDevelopment)
	}
	var createdBranch bool
	if !main.Succeeded() {
		ExecCmd("git", "symbolic-ref", "HEAD", fmt.Sprintf("refs/heads/%s", masterName))
		ExecCmd("git", "commit", "--allow-empty", "--quiet", "-m", "Initial commit")
		createdBranch = true
	}
	if !main.HasLocalBranch(devName) {
		if main.HasRemoteBranch(devName) {
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
			ParsePrefix(context, Key, Value)
		}
		for _, prefix := range DefaultTags {
			value := Get(Key)
			defaultValue := value
			if context.Bool("force") || "" == value {
				var newValue string
				if context.Bool("default") {
					newValue = Value
				} else {
					if "" == value {
						defaultValue = Value
					}
					newValue = PromptForInput(
						TAG_NAME_VALIDATION,
						fmt.Sprintf("Prefix for %s tags? [%s]", Value, defaultValue),
						defaultValue,
					)
				}
				logrus.Trace(Key, newValue)
				Write(Key, newValue)
			}
		}
	}
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Debug("CommandInitAction")
	return InitProcedural(context)
}
