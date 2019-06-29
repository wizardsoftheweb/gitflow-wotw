package main

const (
	UnsetOptionValue = ""
)

type ConfigOptionArgs struct {
	Section    string
	Subsection string
	Key        string
	Value      string
}

const (
	GitflowBranchMasterOption      = ConfigOptionArgs{"gitflow", "branches", "master", "master"}
	GitflowBranchDevelopmentOption = ConfigOptionArgs{"gitflow", "branches", "development", "dev"}
	GitflowPrefixFeatureOption     = ConfigOptionArgs{"gitflow", "prefix", "feature", "feature"}
	GitflowPrefixReleaseOption     = ConfigOptionArgs{"gitflow", "prefix", "release", "release"}
	GitflowPrefixHotfixOption      = ConfigOptionArgs{"gitflow", "prefix", "hotfix", "hotfix"}
	GitflowPrefixSupportOption     = ConfigOptionArgs{"gitflow", "prefix", "support", "support"}
	GitflowPrefixVersiontagOption  = ConfigOptionArgs{"gitflow", "prefix", "versiontag", "v"}
)
