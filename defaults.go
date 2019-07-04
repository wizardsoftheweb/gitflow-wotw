package main

const (
	MASTER_BRANCH_KEY  = "gitflow.branch.master"
	DEV_BRANCH_KEY     = "gitflow.branches.development"
	FEATURE_PREFIX_KEY = "gitflow.prefix.feature"
	RELEASE_PREFIX_KEY = "gitflow.prefix.release"
	HOTFIX_PREFIX_KEY  = "gitflow.prefix.hotfix"
	SUPPORT_PREFIX_KEY = "gitflow.prefix.support"
	VERSIONTAG_KEY     = "gitflow.prefix.versiontag"
)

type Option struct {
	Key   string
	Value string
}

var (
	DefaultBranchMaster      = Option{MASTER_BRANCH_KEY, "master"}
	DefaultBranchDevelopment = Option{DEV_BRANCH_KEY, "dev"}
	DefaultPrefixFeature     = Option{FEATURE_PREFIX_KEY, "feature"}
	DefaultPrefixRelease     = Option{RELEASE_PREFIX_KEY, "release"}
	DefaultPrefixHotfix      = Option{HOTFIX_PREFIX_KEY, "hotfix"}
	DefaultPrefixSupport     = Option{SUPPORT_PREFIX_KEY, "support"}
	DefaultPrefixVersiontag  = Option{VERSIONTAG_KEY, "v"}
)

var (
	DefaultMasterSuggestions = []string{
		"master", "prod", "production", "main",
	}
	DefaultDevSuggestions = []string{
		"develop", "dev", "development", "master",
	}
	DefaultBranches = []Option{
		DefaultBranchDevelopment,
		DefaultBranchMaster,
	}
	DefaultPrefixes = []Option{
		DefaultPrefixFeature,
		DefaultPrefixHotfix,
		DefaultPrefixRelease,
		DefaultPrefixSupport,
	}
	DefaultTags = []Option{
		DefaultPrefixVersiontag,
	}
)
