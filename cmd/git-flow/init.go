package gitflow

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wat")
	},
}

var Defaults bool
var Force bool

func init() {
	GitFlowCmd.AddCommand(InitCmd)
	InitCmd.Flags().BoolVarP(
		&Defaults,
		"defaults",
		"d",
		false,
		"Use defaults everywhere",
	)
	InitCmd.Flags().BoolVarP(
		&Force,
		"force",
		"f",
		false,
		"Force reinitialization",
	)
}
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
func SharedPrep(cmd *cobra.Command, args []string, branchName string) error {
	if cmd.Flag("force") || !IsBranchConfigured(branchName) {
		var checkExistence bool
		var suggestion string
		if 0 == len(LocalBranches()) {
			checkExistence = false
			suggestion = GetWithDefault(GetKeyFromRef(branchName), branchName)
		} else {
			checkExistence = true
			suggestion = PickGoodSuggestion(branchName)
			if "" == suggestion {
				checkExistence = false
				suggestion = GetWithDefault(GetKeyFromRef(branchName), GetValueFromRef(branchName))
			}
		}
		newValue := PromptForInput(REF_NAME_VALIDATION, PromptMessageFromBranch(branchName, suggestion), suggestion)
		Write(GetKeyFromRef(branchName), newValue)
		if checkExistence {
			CheckExistence(branchName, newValue)
		}
	}
	return nil
}

func CheckExistence(branchName string, newName string) {
	if "master" == branchName {
		if !HasLocalBranch(newName) {
			if HasRemoteBranch(newName) {
				ExecCmd("git", "branch", newName, fmt.Sprintf("origin/%s", newName))
			} else {
				logrus.Warn("Master", HasLocalBranch(newName))
				logrus.Warn(fmt.Sprintf("'%s'", newName))
				logrus.Error(fmt.Sprintf("Branch '%s' does not exist", newName))
				logrus.Fatal(ErrProdDoesntExist)
			}
		}
		return
	} else {
		if !HasLocalBranch(newName) {
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

func ParsePrefix(cmd *cobra.Command, args []string, prefixKey string, defaultValue string) error {
	if !cmd.Flag("force").Value.Bool {
		return nil
	} else {
		var newValue string
		if !cmd.Flag("default").Value.Bool {
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

func InitProcedural(cmd *cobra.Command, args []string) error {
	if !Succeeded() {
		GitInit()
	} else {
		if Succeeded() {
			IsWorkingTreeClean()
		} else {
			logrus.Fatal(ErrHeadlessRepo)
		}
	}
	logrus.Trace("Repo has be identified")
	if IsGitFlowInitialized() && !context.Bool("force") {
		logrus.Fatal(ErrAlreadyInitialized)
	}
	if cmd.Flag("default").Value.Bool {
		//fmt.Fprint(context.App.Writer, "Using default branch names")
	}
	for _, branchName := range []string{"master", "dev"} {
		SharedPrep(cmd, args, branchName)
	}
	devName := Get(DEV_BRANCH_KEY)
	masterName := Get(MASTER_BRANCH_KEY)
	if devName == masterName {
		logrus.Fatal(ErrProductionMustDifferFromDevelopment)
	}
	var createdBranch bool
	if !Succeeded() {
		ExecCmd("git", "symbolic-ref", "HEAD", fmt.Sprintf("refs/heads/%s", masterName))
		ExecCmd("git", "commit", "--allow-empty", "--quiet", "-m", "Initial commit")
		createdBranch = true
	}
	if !HasLocalBranch(devName) {
		if HasRemoteBranch(devName) {
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
		//fmt.Fprint(context.App.Writer, "Some prefixes need to be configured")
		for _, prefix := range DefaultPrefixes {
			ParsePrefix(cmd, args, Key, Value)
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
