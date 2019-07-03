package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

type GitConfigOption struct {
	Key   string
	Value string
}

type GitConfigSection struct {
	Heading    string
	Subheading string
	Options    []GitConfigOption
}

type GitConfig struct {
	Sections []GitConfigSection
}

type ConfigFileHandler struct {
	dotGitDir   FileSystemObject
	configFile  FileSystemObject
	rawContents string
}

var (
	GitConfigBlockPattern   = regexp.MustCompile(`(?m)(?:^\[.*?$\s*)(^\s+.*?$\s?)+`)
	GitConfigSectionPattern = regexp.MustCompile(`\[\s*(?P<heading>.*?)(\s+["'](?P<subheading>.*?)["'])?\s*\]`)
	GitConfigOptionPattern  = regexp.MustCompile(`\s+(?P<key>.*?)\s*=\s*(?P<value>.*?)\s*`)
)

func (handler *ConfigFileHandler) createIfDoesNotExist() error {
	fmt.Println("cool")
	return nil
}

func (handler *ConfigFileHandler) exist() bool {
	return handler.configFile.exists()
}

func (handler *ConfigFileHandler) loadConfig() error {
	raw_data, err := ioutil.ReadFile(handler.configFile.String())
	if nil != err {
		log.Fatal(err)
	}
	handler.rawContents = string(raw_data)
	return nil
}

func (handler *ConfigFileHandler) parseOptionConfig(raw_config string) ([]GitConfigOption, error) {
	options := []GitConfigOption{}
	for _, match := range GitConfigOptionPattern.FindAllString(raw_config, -1) {
		result := map[string]string{}
		for index, name := range match {
			result[GitConfigOptionPattern.SubexpNames()[index]] = string(name)
		}
		key, ok := result["key"]
		if !ok {
			log.Fatal(ok)
		}
		value, ok := result["value"]
		if !ok {
			log.Fatal(ok)
		}
		options = append(options, GitConfigOption{
			Key:   key,
			Value: value,
		})
	}
	return options, nil
}

func (handler *ConfigFileHandler) parseSectionConfig(raw_config string) (GitConfigSection, error) {
	section := GitConfigSection{}
	for _, match := range GitConfigSectionPattern.FindAllString(raw_config, -1) {
		result := map[string]string{}
		for index, name := range match {
			result[GitConfigSectionPattern.SubexpNames()[index]] = string(name)
		}
		heading, ok := result["heading"]
		if !ok {
			log.Fatal(ok)
		}
		subheading, ok := result["subheading"]
		if !ok {
			log.Fatal(ok)
		}
		section.Heading = heading
		section.Subheading = subheading
	}
	return section, nil
}

func (handler *ConfigFileHandler) parseBlockConfig(raw_config string) (GitConfigSection, error) {
	section, err := handler.parseSectionConfig(raw_config)
	if nil != err {
		log.Fatal(err)
	}
	options, err := handler.parseOptionConfig(raw_config)
	if nil != err {
		log.Fatal(err)
	}
	section.Options = options
	return section, nil
}
func (handler *ConfigFileHandler) parseConfig() ([]GitConfigSection, error) {
	sections := []GitConfigSection{}
	for _, block := range GitConfigBlockPattern.FindAllString(handler.rawContents, -1) {
		section, err := handler.parseBlockConfig(block)
		if nil != err {
			log.Fatal(err)
		}
		sections = append(sections, section)
	}
	return sections, nil
}
