package gitflow

func RevParseGitDir() CommandResponse {
	return ExecCmd("git", "rev-parse", "--git-dir")
}

func RevParseArgs(arguments ...string) CommandResponse {
	return ExecCmd(append([]string{"git", "rev-parse"}, arguments...)...)
}
func RevParseQuietVerifyHead() CommandResponse {
	return ExecCmd("git", "rev-parse", "--quiet", "--verify", "HEAD")
}

func RevParseAbbrevRefHead() CommandResponse {
	return ExecCmd("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func BranchNoColor(remote bool) CommandResponse {
	action := []string{"git", "branch", "--no-color"}
	if remote {
		action = append(action, "-r")
	}
	return ExecCmd(action...)
}

func MergeBase(firstBranch string, secondBranch string) CommandResponse {
	return ExecCmd("git", "merge-base", firstBranch, secondBranch)
}

func GitInit() CommandResponse {
	return ExecCmd("git", "init")
}
