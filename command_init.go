package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

var (
	git = &Git{}
)

var (
	ErrUnableToGitInit = errors.New("Unable to complete git init in the current working directory")
)

var (
	CommandInit = cli.Command{
		Name: "init",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "force, f",
				Usage: "Forces the setting to be reinitialized",
			},
			cli.BoolFlag{
				Name:  "defaults, d",
				Usage: "Applies the defaults without prompting (when available)",
			},
		},
		Action: CommandInitAction,
	}
)

func EnsureRepoIsAvailable(directory string) (Repository, error) {
	var repo Repository
	var err error
	result := execute("git", "rev-parse", "--git-dir")
	if !result.Succeeded() {
		result = execute("git", "init")
		if !result.Succeeded() {
			return repo, ErrUnableToGitInit
		}
	} else {
		revparse := execute("git", "rev-parse", "--quiet", "--verify", "HEAD")
		if !revparse.Succeeded() {
			println("rad")
		}

	}
	parseExitCode(err)
	return repo, nil
}

func CommandInitAction(context *cli.Context) error {
	directory, _ := os.Getwd()
	EnsureRepoIsAvailable(directory)
	// dot_dir, _ := repo.discoverDotDir(FileSystemObject(directory))
	return nil
}
