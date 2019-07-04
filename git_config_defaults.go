package main

type GitflowDefaultOptions struct {
	Section    string
	Subsection string
	Key        string
	Value      string
}

var (
	DefaultGitflowBranchMasterOption      = GitflowDefaultOptions{"gitflow", "branches", "master", "master"}
	DefaultGitflowBranchDevelopmentOption = GitflowDefaultOptions{"gitflow", "branches", "development", "dev"}
	DefaultGitflowPrefixFeatureOption     = GitflowDefaultOptions{"gitflow", "prefix", "feature", "feature"}
	DefaultGitflowPrefixReleaseOption     = GitflowDefaultOptions{"gitflow", "prefix", "release", "release"}
	DefaultGitflowPrefixHotfixOption      = GitflowDefaultOptions{"gitflow", "prefix", "hotfix", "hotfix"}
	DefaultGitflowPrefixSupportOption     = GitflowDefaultOptions{"gitflow", "prefix", "support", "support"}
	DefaultGitflowPrefixVersiontagOption  = GitflowDefaultOptions{"gitflow", "prefix", "versiontag", "v"}
)

var OptionsToInitializeGitflow = []GitflowDefaultOptions{
	DefaultGitflowBranchMasterOption,
	DefaultGitflowBranchDevelopmentOption,
	DefaultGitflowPrefixFeatureOption,
	DefaultGitflowPrefixReleaseOption,
	DefaultGitflowPrefixHotfixOption,
	DefaultGitflowPrefixSupportOption,
	DefaultGitflowPrefixVersiontagOption,
}

var SubbranchOptions = []GitflowDefaultOptions{
	DefaultGitflowPrefixFeatureOption,
	DefaultGitflowPrefixReleaseOption,
	DefaultGitflowPrefixHotfixOption,
	DefaultGitflowPrefixSupportOption,
}
