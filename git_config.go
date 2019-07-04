package main

import (
	"errors"
	"fmt"
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
	args := fmt.Sprintf("%s", arguments)
	name, ok := OptionKeyCache[args]
	if !ok {
		name := strings.Join(arguments, ".")
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
	switch action {
	case GIT_CONFIG_CREATE:
		return config.create(arguments...)
	case GIT_CONFIG_READ:
		return config.read(arguments...)
	case GIT_CONFIG_UPDATE:
		return config.update(arguments...)
	case GIT_CONFIG_DELETE:
		return config.delete(arguments...)
	}
	return "", nil
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
