package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/sirupsen/logrus"
)

type ConfigFileHandler struct {
	ConfigStorageHandler
	configFile FileSystemObject
}

var (
	GitConfigFileBlockPattern   = regexp.MustCompile(`(?m)(?:^\[.*?$\s*)(^\s+.*?$\s?)+`)
	GitConfigFileSectionPattern = regexp.MustCompile(`(?m)^\[\s*(?P<heading>.*?)(\s+["'](?P<subheading>.*?)["'])?\s*\]\s*$`)
	GitConfigFileOptionPattern  = regexp.MustCompile(`(?m)^\s+(?P<key>.*?)\s*=\s*(?P<value>.*)\s*$`)
)

func (handler *ConfigFileHandler) createIfDoesNotExist() error {
	logrus.Trace("createIfDoesNotExist")
	fmt.Println("cool")
	return nil
}

func (handler *ConfigFileHandler) exist() bool {
	logrus.Trace("exist")
	return handler.configFile.exists()
}

func (handler *ConfigFileHandler) loadConfig() error {
	logrus.Trace("loadConfig")
	logrus.Debug(handler.configFile.String())
	raw_data, err := ioutil.ReadFile(handler.configFile.String())
	if nil != err {
		log.Fatal(err)
	}
	handler.rawContents = string(raw_data)
	return nil
}

func (handler *ConfigFileHandler) parseOptionConfig(config *GitConfig, section []string, raw_config string) (*GitConfig, error) {
	logrus.Trace("parseOptionConfig")
	for _, match := range GitConfigFileOptionPattern.FindAllStringSubmatch(raw_config, -1) {
		result := map[string]string{}
		for index, name := range GitConfigFileOptionPattern.SubexpNames() {
			if 0 != index && "" != name {
				result[name] = match[index]
			}
		}
		config.Option(
			GIT_CONFIG_CREATE,
			append(section, []string{result["key"], result["value"]}...)...,
		)
	}
	return config, nil
}

func (handler *ConfigFileHandler) parseSectionConfig(raw_config string) ([]string, error) {
	logrus.Trace("parseSectionConfig")
	section := []string{}
	for _, match := range GitConfigFileSectionPattern.FindAllStringSubmatch(raw_config, -1) {
		result := map[string]string{}
		for index, name := range GitConfigFileSectionPattern.SubexpNames() {
			if 0 != index && "" != name {
				result[name] = string(match[index])
			}
		}
		returnValue := []string{result["heading"]}
		if "" != result["subheading"] {
			return append(returnValue, result["subheading"]), nil
		}
		return returnValue, nil
	}
	return section, nil
}

func (handler *ConfigFileHandler) parseBlockConfig(config *GitConfig, raw_config string) (*GitConfig, error) {
	logrus.Trace("parseBlockConfig")
	section, err := handler.parseSectionConfig(raw_config)
	if nil != err {
		log.Fatal(err)
	}
	update_config, err := handler.parseOptionConfig(config, section, raw_config)
	if nil != err {
		log.Fatal(err)
	}
	return update_config, nil
}
func (handler *ConfigFileHandler) parseConfig() (*GitConfig, error) {
	logrus.Trace("parseConfig")
	config := &GitConfig{
		Options: make(map[string]*GitOption),
	}
	for _, block := range GitConfigFileBlockPattern.FindAllString(handler.rawContents, -1) {
		result, err := handler.parseBlockConfig(config, block)
		if nil != err {
			log.Fatal(err)
		}
		config = result
	}

	return config, nil
}

func (handler *ConfigFileHandler) dumpConfig(config *GitConfig) []string {
	logrus.Trace("dumpConfig")
	lines := []string{}
	currentSection := ""
	oldSection := "qqq"
	for _, key := range config.SortKeys() {
		option, _ := config.Options[key]
		logrus.Debug(fmt.Sprintf("%s: %s", key, option))
		if "" != option.Subsection {
			currentSection = FormatGitConfigSectionFileName(option.Section, option.Subsection)
		} else {
			currentSection = FormatGitConfigSectionFileName(option.Section)
		}
		if oldSection != currentSection {
			lines = append(lines, currentSection)
			oldSection = currentSection
		}
		lines = append(
			lines,
			fmt.Sprintf("\t%s = %s", option.Key, option.Value),
		)
	}
	return lines
}
