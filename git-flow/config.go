package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	format "gopkg.in/src-d/go-git.v4/plumbing/format/config"
)

func GetOptionValue(git_config *format.Config, section string, subsection string, key string) string {
	logrus.Debug("GetOptionValue")
	if format.NoSubsection == subsection {
		return git_config.Section(section).Option(key)
	} else {
		return git_config.Section(section).Subsection(subsection).Option(key)
	}
	return UnsetOptionValue
}

func IsOptionSet(git_config *format.Config, section string, subsection string, key string) bool {
	logrus.Debug("IsOptionSet")
	return UnsetOptionValue != GetOptionValue(git_config, section, subsection, key)
}

func (config_options_args *ConfigOptionArgs) New(section string, subsection string, key string, value string) *ConfigOptionArgs {
	logrus.Debug("New")
	config_options_args.Section = section
	config_options_args.Subsection = subsection
	config_options_args.Key = key
	config_options_args.Value = value
	return config_options_args
}

func (config_options_args *ConfigOptionArgs) isOptionSetInConfig(git_config *format.Config) bool {
	logrus.Debug("isOptionSetInConfig")
	return IsOptionSet(git_config, config_options_args.Section, config_options_args.Subsection, config_options_args.Key)
}

func LoadConfig(repo *git.Repository) (*config.Config, error) {
	logrus.Debug("LoadConfig")
	git_config, err := repo.Config()
	CheckError(err)
	return git_config, nil
}
