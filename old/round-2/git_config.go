package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

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

type OptionKeyStore map[string]string

var (
	OptionKeyCache = make(OptionKeyStore)
)

type GitConfigCrud int

const (
	GIT_CONFIG_CREATE GitConfigCrud = iota
	GIT_CONFIG_READ
	GIT_CONFIG_UPDATE
	GIT_CONFIG_DELETE
)

type GitOption struct {
	Section    string
	Subsection string
	Key        string
	Value      string
}

func NewGitOption(arguments ...string) *GitOption {
	newOption := &GitOption{}
	newOption.Value, arguments = arguments[len(arguments)-1], arguments[:len(arguments)-1]
	newOption.Key, arguments = arguments[len(arguments)-1], arguments[:len(arguments)-1]
	newOption.Section, arguments = arguments[0], arguments[1:]
	if 1 == len(arguments) {
		newOption.Subsection = arguments[0]
	}
	return newOption
}

type GitConfig struct {
	Options map[string]*GitOption
}

func GenerateOptionName(arguments ...string) string {
	logrus.Trace("GenerateOptionName")
	args := fmt.Sprintf("%s", arguments)
	name, ok := OptionKeyCache[args]
	if !ok {
		name = strings.Join(arguments, ".")
		OptionKeyCache[args] = name
	}
	return name
}

func (config *GitConfig) create(arguments ...string) (string, error) {
	key := GenerateOptionName(arguments[:len(arguments)-1]...)
	value := NewGitOption(arguments...)
	config.Options[key] = value
	return "", nil

}

func (config *GitConfig) read(arguments ...string) (string, error) {
	key := GenerateOptionName(arguments...)
	value, ok := config.Options[key]
	if !ok {
		return "", ErrConfigOptionNotFound
	}
	return value.Value, nil
}

func (config *GitConfig) update(arguments ...string) (string, error) {
	key := GenerateOptionName(arguments[:len(arguments)-1]...)
	value := NewGitOption(arguments...)
	config.Options[key] = value
	return "", nil
}

func (config *GitConfig) delete(arguments ...string) (string, error) {
	key := GenerateOptionName(arguments...)
	delete(config.Options, key)
	return "", nil
}

func (config *GitConfig) Option(action GitConfigCrud, arguments ...string) (string, error) {
	logrus.Trace("Option")
	switch action {
	case GIT_CONFIG_CREATE:
		logrus.Debug(fmt.Sprintf("%s: %s", "create", arguments))
		return config.create(arguments...)
	case GIT_CONFIG_READ:
		logrus.Debug(fmt.Sprintf("%s: %s", "read", arguments))
		return config.read(arguments...)
	case GIT_CONFIG_UPDATE:
		logrus.Debug(fmt.Sprintf("%s: %s", "update", arguments))
		return config.update(arguments...)
	case GIT_CONFIG_DELETE:
		logrus.Debug(fmt.Sprintf("%s: %s", "delete", arguments))
		return config.delete(arguments...)
	}
	return "", nil
}

func (config *GitConfig) SortKeys() []string {
	keys := make([]string, len(config.Options))
	index := 0
	for key, _ := range config.Options {
		keys[index] = key
		index++
	}
	sort.Strings(keys)
	return keys
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
