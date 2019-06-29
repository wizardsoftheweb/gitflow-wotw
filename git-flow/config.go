package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	format "gopkg.in/src-d/go-git.v4/plumbing/format/config"
)

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
	GitflowBranchMasterOption      = ConfigOptionArgs{"gitflow", "branches", "master", UnsetOptionValue}
	GitflowBranchDevelopmentOption = ConfigOptionArgs{"gitflow", "branches", "development", UnsetOptionValue}
	GitflowPrefixFeatureOption     = ConfigOptionArgs{"gitflow", "prefix", "feature", UnsetOptionValue}
	GitflowPrefixReleaseOption     = ConfigOptionArgs{"gitflow", "prefix", "release", UnsetOptionValue}
	GitflowPrefixHotfixOption      = ConfigOptionArgs{"gitflow", "prefix", "hotfix", UnsetOptionValue}
	GitflowPrefixSupportOption     = ConfigOptionArgs{"gitflow", "prefix", "support", UnsetOptionValue}
	GitflowPrefixVersiontagOption  = ConfigOptionArgs{"gitflow", "prefix", "versiontag", UnsetOptionValue}
)

func GetOptionValue(git_config *format.Config, section string, subsection string, key string) string {
	if format.NoSubsection == subsection {
		return git_config.Section(section).Option(key)
	} else {
		return git_config.Section(section).Subsection(subsection).Option(key)
	}
	return UnsetOptionValue
}

func IsOptionSet(git_config *format.Config, section string, subsection string, key string) bool {
	return UnsetOptionValue != GetOptionValue(git_config, section, subsection, key)
}

func (config_options_args *ConfigOptionArgs) New(section string, subsection string, key string, value string) *ConfigOptionArgs {
	config_options_args.Section = section
	config_options_args.Subsection = subsection
	config_options_args.Key = key
	config_options_args.Value = value
	return config_options_args
}

func (config_options_args *ConfigOptionArgs) isOptionSetInConfig(git_config *format.Config) bool {
	return IsOptionSet(git_config, config_options_args.Section, config_options_args.Subsection, config_options_args.Key)
}

var (
	NecessaryInitSettings = []ConfigOptionArgs{
		GitflowBranchDevelopmentOption,
		GitflowBranchMasterOption,
		GitflowPrefixFeatureOption,
		GitflowPrefixHotfixOption,
		GitflowPrefixReleaseOption,
		GitflowPrefixSupportOption,
		GitflowPrefixVersiontagOption,
	}
)

func EnsureNecessaryInitOptionsAreSet(git_config *format.Config) bool {
	for _, option_arg_set := range NecessaryInitSettings {
		if !option_arg_set.isOptionSetInConfig(git_config) {
			return false
		}
	}
	return true

}

func LoadConfig(repo *git.Repository) (*config.Config, error) {
	git_config, err := repo.Config()
	CheckError(err)
	return git_config, nil
}
