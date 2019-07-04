package main

type Option struct {
	Key   string
	Value string
}

var (
	DefaultBranchMaster      = Option{"gitflow.branches.master", "master"}
	DefaultBranchDevelopment = Option{"gitflow.branches.development", "dev"}
	DefaultPrefixFeature     = Option{"gitflow.prefix.feature", "feature"}
	DefaultPrefixRelease     = Option{"gitflow.prefix.release", "release"}
	DefaultPrefixHotfix      = Option{"gitflow.prefix.hotfix", "hotfix"}
	DefaultPrefixSupport     = Option{"gitflow.prefix.support", "support"}
	DefaultPrefixVersiontag  = Option{"gitflow.prefix.versiontag", "v"}
)

var (
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
