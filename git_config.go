package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
)

var (
	ErrConfigOptionNotFound       = errors.New("That config option has not been set")
	ErrConfigSectionNotFound      = errors.New("That config section has not been built")
	ErrUnableToProcessCrudRequest = errors.New("Unablet to process your CRUD request")
)

type ConfigStorageHandler struct {
	dotGitDir   FileSystemObject
	rawContents string
}

type GitConfigCrud int

const (
	GIT_CONFIG_CREATE GitConfigCrud = iota
	GIT_CONFIG_READ
	GIT_CONFIG_UPDATE
	GIT_CONFIG_DELETE
)

type GitConfigOptions map[string]string

type GitConfigSection struct {
	Heading    string
	Subheading string
	Options    GitConfigOptions
}

type SectionHeaderFormats int

const (
	HEADER_WITH_SUBHEADING SectionHeaderFormats = iota
	HEADER_WITHOUT
)

var (
	GitConfigSectionFileHeaderFormats = []string{
		"[%s \"%s\"]",
		"[%s]",
	}
	GitConfigSectionEnvironmentFormats = []string{
		"%s.%s",
		"%s",
	}
)

func FormatGitConfigSectionFileName(components ...string) string {
	logrus.Trace("FormatGitConfigSectionFileName")
	if "" == components[0] {
		return ""
	}
	if 1 < len(components) && "" != components[1] {
		return fmt.Sprintf(GitConfigSectionFileHeaderFormats[HEADER_WITH_SUBHEADING], components[0], components[1])
	}
	return fmt.Sprintf(GitConfigSectionFileHeaderFormats[HEADER_WITHOUT], components[0])
}
func FormatGitConfigSectionEnvironmentName(components ...string) string {
	logrus.Trace("FormatGitConfigSectionEnvironmentName")
	if "" == components[0] {
		return ""
	}
	if 1 < len(components) && "" != components[1] {
		return fmt.Sprintf(GitConfigSectionEnvironmentFormats[HEADER_WITH_SUBHEADING], components[0], components[1])
	}
	return fmt.Sprintf(GitConfigSectionEnvironmentFormats[HEADER_WITHOUT], components[0])
}

func (section *GitConfigSection) FileHeader() string {
	logrus.Trace("FileHeader")
	if 0 < len(section.Subheading) {
		return FormatGitConfigSectionFileName(section.Heading, section.Subheading)
	}
	return FormatGitConfigSectionFileName(section.Heading)
}

func (section *GitConfigSection) EnvironmentHeader() string {
	logrus.Trace("EnvironmentHeader")
	if 0 < len(section.Subheading) {
		return FormatGitConfigSectionEnvironmentName(section.Heading, section.Subheading)
	}
	return FormatGitConfigSectionEnvironmentName(section.Heading)
}

func (section *GitConfigSection) create(key string, value string) error {
	logrus.Trace("create")
	section.Options[key] = value
	return nil
}

func (section *GitConfigSection) read(key string) (string, error) {
	logrus.Trace(fmt.Sprintf("read %s.%s", section.EnvironmentHeader(), key))
	value, ok := section.Options[key]
	if !ok {
		return "", ErrConfigOptionNotFound
	}
	return value, nil
}

func (section *GitConfigSection) update(key string, value string) error {
	logrus.Trace("update")
	section.Options[key] = value
	return nil
}

func (section *GitConfigSection) delete(key string) error {
	logrus.Trace("delete")
	delete(section.Options, key)
	return nil
}

type GitConfig struct {
	Sections map[string]GitConfigSection
}

func (config *GitConfig) create(new_config GitConfigSection) error {
	logrus.Trace("create")
	config.Sections[new_config.FileHeader()] = new_config
	return nil
}

func (config *GitConfig) read(components ...string) (GitConfigSection, error) {
	logrus.Trace("read")
	value, ok := config.Sections[FormatGitConfigSectionFileName(components...)]
	if !ok {
		value, ok = config.Sections[components[0]]
		if !ok {
			return GitConfigSection{}, ErrConfigSectionNotFound
		}
	}
	return value, nil
}

func (config *GitConfig) update(key string, value GitConfigSection) error {
	logrus.Trace("update")
	config.Sections[key] = value
	return nil
}

func (config *GitConfig) delete(key string) error {
	logrus.Trace("delete")
	delete(config.Sections, key)
	return nil
}

func (config *GitConfig) Option(action GitConfigCrud, components ...string) (string, error) {
	logrus.Trace("Option")
	var section_name string
	var key, value int
	if 4 == len(components) {
		section_name = FormatGitConfigSectionFileName(components[0], components[1])
		key = 2
		value = 3
		_, err := config.read(section_name)
		if nil != err {
			config.create(GitConfigSection{
				Heading:    components[0],
				Subheading: components[1],
				Options:    make(map[string]string),
			})
		}
	} else {
		section_name = FormatGitConfigSectionFileName(components[0])
		key = 1
		value = 2
		_, err := config.read(section_name)
		if nil != err {
			config.create(GitConfigSection{
				Heading: components[0],
				Options: make(map[string]string),
			})
		}
	}
	section, err := config.read(section_name)
	if nil != err {
		logrus.Warning(err)
		log.Fatal(err)
	}

	switch {
	case GIT_CONFIG_CREATE == action && 3 <= len(components):
		return "", section.create(
			components[key],
			components[value],
		)
	case GIT_CONFIG_READ == action:
		return section.read(
			components[key],
		)
	case GIT_CONFIG_UPDATE == action && 3 <= len(components):
		return "", section.update(
			components[key],
			components[value],
		)
	case GIT_CONFIG_DELETE == action:
		return "", section.delete(
			components[key],
		)
	}
	return "", ErrUnableToProcessCrudRequest
}
