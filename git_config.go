package main

type GitConfigHelper struct {
}

var (
	GitConfig = &GitConfigHelper{}
)

func (g *GitConfigHelper) Get(key string) string {
	result := ExecCmd("git", "config", "--get", key)
	return result.result
}

func (g *GitConfigHelper) Write(key string, value string) CommandResponse {
	result := ExecCmd("git", "config", "--local", key, value)
	return result
}
