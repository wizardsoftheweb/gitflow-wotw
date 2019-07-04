package main

func RevParseGitDir() CommandResponse {
	return ExecCmd("git", "rev-parse", "--git-dir")
}

func RevParseArgs(arguments ...string) CommandResponse {
	return ExecCmd(append([]string{"git", "rev-parse"}, arguments...)...)
}
func RevParseQuietVerifyHead() CommandResponse {
	return ExecCmd("git", "rev-parse", "--quiet", "--verify", "HEAD")
}

func BranchNoColor(remote bool) CommandResponse {
	action := []string{"git", "branch", "--no-color"}
	if remote {
		action = append(action, "-r")
	}
	return ExecCmd(action...)
}
