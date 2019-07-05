package gitflow

const (
	MasterBranchKey  = "gitflow.branch.master"
	DevBranchKey     = "gitflow.branch.develop"
	FeaturePrefixKey = "gitflow.prefix.feature"
	ReleasePrefixKey = "gitflow.prefix.release"
	HotfixPrefixKey  = "gitflow.prefix.hotfix"
	SupportPrefixKey = "gitflow.prefix.support"
	VersiontagKey    = "gitflow.prefix.versiontag"
)

type Option struct {
	Key   string
	Value string
}

var (
	DefaultBranchMaster      = Option{MasterBranchKey, "master"}
	DefaultBranchDevelopment = Option{DevBranchKey, "dev"}
	DefaultPrefixFeature     = Option{FeaturePrefixKey, "feature"}
	DefaultPrefixRelease     = Option{ReleasePrefixKey, "release"}
	DefaultPrefixHotfix      = Option{HotfixPrefixKey, "hotfix"}
	DefaultPrefixSupport     = Option{SupportPrefixKey, "support"}
	DefaultPrefixVersiontag  = Option{VersiontagKey, "v"}
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
