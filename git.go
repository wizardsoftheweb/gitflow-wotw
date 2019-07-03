package main

type Git struct {
}

type GitCommandOptions struct{}

func (git *Git) execute(command ...string) CommandResponse {
	stdout, stderr, err := RunCommand(command)
	return CommandResponse{
		stdout:   stdout,
		stderr:   stderr,
		exitCode: parseExitCode(err),
		exitErr:  err,
	}
}

type RevParseOptions struct {
	GitCommandOptions
	GitDir bool
}

func (git *Git) RevParse(options RevParseOptions) CommandResponse {
	command := []string{
		"git",
		"rev-parse",
	}
	if options.GitDir {
		command = append(command, "--git-dir")
	}
	return git.execute(command)
}
