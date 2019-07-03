package main

func IsGitflowInitialized(config GitConfig) bool {
	for _, option := range OptionsToInitializeGitflow {
		value, err := config.Option(GIT_CONFIG_READ, option.Section, option.Subsection, option.Key)
		if nil != err || "" == value {
			return false
		}
	}
	return true
}
