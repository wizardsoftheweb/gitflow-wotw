package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

type ConfigFileHandler struct {
	dotGitDir   FileSystemObject
	configFile  FileSystemObject
	rawContents string
}

var (
	GitConfigBlockPattern   = regexp.MustCompile(`(?m)(?:^\[.*?$\s*)(^\s+.*?$\s?)+`)
	GitConfigSectionPattern = regexp.MustCompile(`(?m)^\[\s*(?P<heading>.*?)(\s+["'](?P<subheading>.*?)["'])?\s*\]\s*$`)
	GitConfigOptionPattern  = regexp.MustCompile(`(?m)^\s+(?P<key>.*?)\s*=\s*(?P<value>.*)\s*$`)
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

func (handler *ConfigFileHandler) parseOptionConfig(raw_config string) (GitConfigOptions, error) {
	options := make(map[string]string)
	for _, match := range GitConfigOptionPattern.FindAllStringSubmatch(raw_config, -1) {
		result := map[string]string{}
		for index, name := range GitConfigOptionPattern.SubexpNames() {
			if 0 != index && "" != name {
				result[name] = match[index]
			}
		}
		key := result["key"]
		value := result["value"]
		options[key] = value
		fmt.Println(key, value)
	}
	return options, nil
}

func (handler *ConfigFileHandler) parseSectionConfig(raw_config string) (GitConfigSection, error) {
	section := GitConfigSection{}
	for _, match := range GitConfigSectionPattern.FindAllStringSubmatch(raw_config, -1) {
		result := map[string]string{}
		for index, name := range GitConfigSectionPattern.SubexpNames() {
			if 0 != index && "" != name {
				result[name] = string(match[index])
			}
		}
		section.Heading = result["heading"]
		section.Subheading = result["subheading"]
	}
	return section, nil
}

func (handler *ConfigFileHandler) parseBlockConfig(raw_config string) (GitConfigSection, error) {
	fmt.Println("-----------")
	fmt.Println(raw_config)
	fmt.Println("-----------")
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
func (handler *ConfigFileHandler) parseConfig() (GitConfig, error) {
	sections := make(map[string]GitConfigSection)
	for _, block := range GitConfigBlockPattern.FindAllString(handler.rawContents, -1) {
		section, err := handler.parseBlockConfig(block)
		if nil != err {
			log.Fatal(err)
		}
		sections[section.Name()] = section
	}

	return GitConfig{
		Sections: sections,
	}, nil
}

func (handler *ConfigFileHandler) dumpOption(options GitConfigOptions) []string {
	lines := []string{}
	for key, value := range options {
		lines = append(
			lines,
			fmt.Sprintf("\t%s = %s", key, value),
		)
	}
	return lines
}

func (handler *ConfigFileHandler) dumpSection(section GitConfigSection) []string {
	lines := []string{
		section.Name(),
	}
	return append(lines, handler.dumpOption(section.Options)...)
}

func (handler *ConfigFileHandler) dumpConfig(config GitConfig) []string {
	lines := []string{}
	for _, section := range config.Sections {
		lines = append(lines, handler.dumpSection(section)...)
	}
	return lines
}
