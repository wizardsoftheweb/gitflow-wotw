package gitflow

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:              "init",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wat")
	},
}

var Defaults bool
var Force bool

func init() {
	PackageCmd.AddCommand(InitCmd)
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
		return MasterBranchKey
	}
	return DevBranchKey
}

func GetValueFromRef(branch string) string {
	switch branch {
	case "master":
		return "master"
	}
	return "dev"
}

func SharedPrep(cmd *cobra.Command, args []string, branchName string) error {
	if Force || !IsBranchConfigured(branchName) {
		var checkExistence bool
		var suggestion string
		if 0 == len(Repo.LocalBranches()) {
			checkExistence = false
			suggestion = GitConfig.GetWithDefault(GetKeyFromRef(branchName), branchName)
		} else {
			checkExistence = true
			suggestion = Repo.PickGoodSuggestion(branchName)
			if "" == suggestion {
				checkExistence = false
				suggestion = GitConfig.GetWithDefault(GetKeyFromRef(branchName), GetValueFromRef(branchName))
			}
		}
		newValue := PromptForInput(REF_NAME_VALIDATION, PromptMessageFromBranch(branchName, suggestion), suggestion)
		GitConfig.Write(GetKeyFromRef(branchName), newValue)
		if checkExistence {
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

func ParsePrefix(prefixKey string, defaultValue string) error {
	if !Force {
		return nil
	} else {
		var newValue string
		if Defaults {
			newValue = defaultValue
		} else {
			value := GitConfig.GetWithDefault(prefixKey, defaultValue)
			newValue = PromptForInput(
				PrefixNameValidation,
				fmt.Sprintf("Prefix for %s branches? [%s]", filepath.Base(prefixKey), defaultValue),
				value,
			)
		}
		GitConfig.Write(prefixKey, newValue)
	}
	return nil
}

func InitProcedural(cmd *cobra.Command, args []string) error {
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
	if IsGitFlowInitialized() && !Force {
		logrus.Fatal(ErrAlreadyInitialized)
	}
	if Defaults {
		fmt.Println("Using default branch names")
	}
	for _, branchName := range []string{"master", "dev"} {
		err := SharedPrep(cmd, args, branchName)
		CheckErr(err)
	}
	devName := GitConfig.Get(DevBranchKey)
	masterName := GitConfig.Get(MasterBranchKey)
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

	if Force || !ArePrefixesConfigured() {
		fmt.Println("Some prefixes need to be configured")
		for _, prefix := range DefaultPrefixes {
			err := ParsePrefix(cmd, args, prefix.Key, prefix.Value)
			CheckErr(err)
		}
		for _, prefix := range DefaultTags {
			value := GitConfig.Get(prefix.Key)
			defaultValue := value
			if Force || "" == value {
				var newValue string
				if Defaults {
					newValue = prefix.Value
				} else {
					if "" == value {
						defaultValue = prefix.Value
					}
					newValue = PromptForInput(
						TagNameValidation,
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
