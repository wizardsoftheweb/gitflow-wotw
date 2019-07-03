package main

import (
	"os"

	"github.com/urfave/cli"
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

func EnsureRepoIsAvailable(directory string) Repository {
	var err error
	if !DidCommandSucceed([]string{"git", "rev-parse", "--git-dir"}) {
		_, _, err = executeCommand([]string{"git", "init"})
	} else {
		println("gnarly")
	}
	repo := Repository{}
	return repo
}

func CommandInitAction(context *cli.Context) error {
	directory, _ := os.Getwd()
	EnsureRepoIsAvailable(directory)
	// dot_dir, _ := repo.discoverDotDir(FileSystemObject(directory))
	return nil
}
