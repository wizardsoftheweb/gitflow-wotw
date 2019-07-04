package gitflow

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func IsWorkingTreeClean() bool {
	result := ExecCmd("git", "diff", "--no-ext-diff", "--ignore-submodules", "--quiet", "--exit-code")
	if !main.Succeeded() {
		logrus.Fatal(ErrUnstagedChanges)
	}
	result = ExecCmd("git", "diff-index", "--cached", "--quiet", "--ignore-submodules", "HEAD", "--")
	if !main.Succeeded() {
		logrus.Fatal(ErrIndexUncommitted)
	}
	return true
}

func IsBranchConfigured(name string) bool {
	branchName := Get(fmt.Sprintf("gitflow.branch.%s", name))
	logrus.Trace(branchName)
	return "" != branchName && main.HasLocalBranch(branchName)
}

func IsMasterConfigured() bool {
	return IsBranchConfigured("master")
}

func IsDevConfigured() bool {
	return IsBranchConfigured("dev")
}

func AreMasterAndDevTheSameValue() bool {
	masterName := Get("gitflow.branch.master")
	devName := Get("gitflow.branch.dev")
	return "" != masterName && "" != devName && masterName != devName
}

func ArePrefixesConfigured() bool {
	for _, option := range DefaultPrefixes {
		result := Get(Key)
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
