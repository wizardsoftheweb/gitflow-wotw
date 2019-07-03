package main

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrConfigOptionNotFound       = errors.New("That config option has not been set")
	ErrConfigSectionNotFound      = errors.New("That config section has not been built")
	ErrUnableToProcessCrudRequest = errors.New("Unablet to process your CRUD request")
)

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

func FormatGitConfigSectionName(components ...string) string {
	if "" == components[0] {
		return ""
	}
	if 1 < len(components) && "" != components[1] {
		return fmt.Sprintf("[%s \"%s\"]", components[0], components[1])
	}
	return fmt.Sprintf("[%s]", components[0])
}

func (section *GitConfigSection) Name() string {
	if 0 < len(section.Subheading) {
		return FormatGitConfigSectionName(section.Heading, section.Subheading)
	}
	return FormatGitConfigSectionName(section.Heading)
}

func (section *GitConfigSection) create(key string, value string) error {
	section.Options[key] = value
	return nil
}

func (section *GitConfigSection) read(key string) (string, error) {
	value, ok := section.Options[key]
	if !ok {
		return "", ErrConfigOptionNotFound
	}
	return value, nil
}

func (section *GitConfigSection) update(key string, value string) error {
	section.Options[key] = value
	return nil
}

func (section *GitConfigSection) delete(key string) error {
	delete(section.Options, key)
	return nil
}

type GitConfig struct {
	Sections map[string]GitConfigSection
}

func (config *GitConfig) create(new_config GitConfigSection) error {
	config.Sections[new_config.Name()] = new_config
	return nil
}

func (config *GitConfig) read(components ...string) (GitConfigSection, error) {
	value, ok := config.Sections[FormatGitConfigSectionName(components...)]
	if !ok {
		value, ok = config.Sections[components[0]]
		if !ok {
			return GitConfigSection{}, ErrConfigSectionNotFound
		}
	}
	return value, nil
}

func (config *GitConfig) update(key string, value GitConfigSection) error {
	config.Sections[key] = value
	return nil
}

func (config *GitConfig) delete(key string) error {
	delete(config.Sections, key)
	return nil
}

func (config *GitConfig) Option(action GitConfigCrud, components ...string) (string, error) {
	var section_name string
	var key, value int
	if 4 == len(components) {
		section_name = FormatGitConfigSectionName(components[0], components[1])
		key = 2
		value = 3
		_, err := config.read(section_name)
		if nil != err {
			config.create(GitConfigSection{
				Heading:    components[0],
				Subheading: components[1],
			})
		}
	} else {
		section_name = FormatGitConfigSectionName(components[0])
		key = 1
		value = 2
		_, err := config.read(section_name)
		if nil != err {
			config.create(GitConfigSection{
				Heading: components[0],
			})
		}
	}
	section, err := config.read(section_name)
	if nil != err {
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
