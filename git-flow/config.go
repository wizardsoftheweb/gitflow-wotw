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
		{"gitflow", "branches", "master", UnsetOptionValue},
		{"gitflow", "branches", "development", UnsetOptionValue},
		{"gitflow", "prefix", "feature", UnsetOptionValue},
		{"gitflow", "prefix", "release", UnsetOptionValue},
		{"gitflow", "prefix", "hotfix", UnsetOptionValue},
		{"gitflow", "prefix", "support", UnsetOptionValue},
		{"gitflow", "prefix", "versiontag", UnsetOptionValue},
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
