package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	GitConfigEnvironmentLinePattern = regexp.MustCompile(`(?m)^\s*(?P<heading>.*?)\.(((?P<subheading>.*)\.)?(?P<key>.*?))\s*=\s*(?P<value>.*)\s*$`)
)

type ConfigEnvironmentHandler struct {
	ConfigStorageHandler
}

func (handler *ConfigEnvironmentHandler) loadConfig() error {
	logrus.Trace("loadConfig")
	command := exec.Command("git", "config", "--local", "--list")
	output, err := command.CombinedOutput()
	if nil != err {
		log.Fatal(err)
	}
	handler.rawContents = string(output)
	return nil
}

func (handler *ConfigEnvironmentHandler) parseConfig() (GitConfig, error) {
	logrus.Trace("parseConfig")
	config := GitConfig{
		Sections: make(map[string]GitConfigSection),
	}
	for _, match := range GitConfigEnvironmentLinePattern.FindAllStringSubmatch(handler.rawContents, -1) {
		result := map[string]string{}
		for index, name := range GitConfigEnvironmentLinePattern.SubexpNames() {
			if 0 != index && "" != name {
				result[name] = match[index]
			}
		}
		section := result["heading"]
		key := result["key"]
		value := result["value"]
		if 4 == len(result) {
			subsection := result["subheading"]
			config.Option(GIT_CONFIG_CREATE, section, subsection, key, value)
			continue
		}
		config.Option(GIT_CONFIG_CREATE, section, key, value)
	}
	return config, nil
}

func (handler *ConfigEnvironmentHandler) dumpConfig(config GitConfig) []string {
	logrus.Trace("dumpConfig")
	lines := []string{}
	for _, section := range config.Sections {
		for key, value := range section.Options {
			lines = append(
				lines,
				fmt.Sprintf(
					"%s.%s=%s",
					section.EnvironmentHeader(),
					key,
					value,
				),
			)
		}
	}
	return lines
}
