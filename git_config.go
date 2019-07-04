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
