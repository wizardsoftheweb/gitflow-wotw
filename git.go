package main

type Git struct {
}

type GitOptions struct{}

func (git *Git) execute(command ...string) CommandResponse {
	if "git" != command[0] {
		command = append([]string{"git"}, command...)
	}
	stdout, stderr, err := RunCommand(command)
	return CommandResponse{
		stdout:   stdout,
		stderr:   stderr,
		exitCode: parseExitCode(err),
		exitErr:  err,
	}
}

func (git *Git) Init() CommandResponse {
	return git.execute("init")
}

type RevParseOptions struct {
	GitOptions
	GitDir bool
	Verify bool
	Quiet  bool
}

func (git *Git) RevParse(options RevParseOptions, arguments ...string) CommandResponse {
	command := []string{
		"rev-parse",
	}
	if options.GitDir {
		command = append(command, "--git-dir")
	}
	if options.Verify {
		command = append(command, "--verify")
	}
	if options.Quiet {
		command = append(command, "--quiet")
	}
	return git.execute(append(command, arguments...)...)
}
