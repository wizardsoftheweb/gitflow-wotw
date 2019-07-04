package main

import "github.com/sirupsen/logrus"

func IsGitflowInitialized(repo *Repository) bool {
	logrus.Trace("IsGitflowInitialized")
	for _, option := range OptionsToInitializeGitflow {
		value, err := repo.config.Option(GIT_CONFIG_READ, option.Section, option.Subsection, option.Key)
		if nil != err || "" == value {
			return false
		}
	}
	if repo.HasBranchBeenConfigured("master") && repo.HasBranchBeenConfigured("dev") {
		return true
	}
	return false
}
