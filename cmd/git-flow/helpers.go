package gitflow

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func CheckErr(err error) {
	if nil != err {
		logrus.Fatal(err)
	}
}

func IsWorkingTreeClean() bool {
	if !ExecCmd("git", "diff", "--no-ext-diff", "--ignore-submodules", "--quiet", "--exit-code").Succeeded() {
		logrus.Fatal(ErrUnstagedChanges)
	}
	if !ExecCmd("git", "diff-index", "--cached", "--quiet", "--ignore-submodules", "HEAD", "--").Succeeded() {
		logrus.Fatal(ErrIndexUncommitted)
	}
	return true
}

func IsBranchConfigured(name string) bool {
	branchName := GitConfig.Get(fmt.Sprintf("gitflow.branch.%s", name))
	logrus.Trace(branchName)
	return "" != branchName && Repo.HasLocalBranch(branchName)
}

func IsMasterConfigured() bool {
	return IsBranchConfigured("master")
}

func IsDevConfigured() bool {
	return IsBranchConfigured("dev")
}

func AreMasterAndDevTheSameValue() bool {
	masterName := GitConfig.Get("gitflow.branch.master")
	devName := GitConfig.Get("gitflow.branch.dev")
	return "" != masterName && "" != devName && masterName != devName
}

func ArePrefixesConfigured() bool {
	for _, option := range DefaultPrefixes {
		result := GitConfig.Get(option.Key)
		if "" == result {
			return false
		} else {
			logrus.Trace(result)
		}
	}
	return true
}

func IsGitFlowInitialized() bool {
	return IsMasterConfigured() &&
		IsDevConfigured() &&
		AreMasterAndDevTheSameValue() &&
		ArePrefixesConfigured()
}

func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
