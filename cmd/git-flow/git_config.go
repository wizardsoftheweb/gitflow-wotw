package gitflow

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type GitConfigHelper struct {
}

var (
	GitConfig = &GitConfigHelper{}
)

func (g *GitConfigHelper) Get(key string) string {
	logrus.Trace(fmt.Sprintf("Get: %s", key))
	result := ExecCmd("git", "config", "--get", key)
	return strings.TrimSpace(result.result)
}

func (g *GitConfigHelper) GetWithDefault(key string, defaultValue string) string {
	logrus.Trace(fmt.Sprintf("Get: %s", key))
	result := ExecCmd("git", "config", "--get", key)
	parsed := strings.TrimSpace(result.result)
	if "" == parsed {
		return defaultValue
	}
	return parsed
}

func (g *GitConfigHelper) Write(key string, value string) CommandResponse {
	logrus.Trace(fmt.Sprintf("Write: %s %s", key, value))
	result := ExecCmd("git", "config", key, strings.TrimSpace(value))
	return result
}
