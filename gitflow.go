package main

import "github.com/sirupsen/logrus"

func IsGitflowInitialized(config GitConfig) bool {
	logrus.Trace("IsGitflowInitialized")
	for _, option := range OptionsToInitializeGitflow {
		value, err := config.Option(GIT_CONFIG_READ, option.Section, option.Subsection, option.Key)
		if nil != err || "" == value {
			return false
		}
	}
	return true
}
