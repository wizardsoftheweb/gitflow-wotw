package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

var (
	ErrUnableToGitInit = errors.New("Unable to complete git init in the current working directory")
	ErrHeadlessRepo    = errors.New("Unable to initialize in a bare repo")
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
		headless := IsRepoHeadless()
		if headless {
			return repo, ErrHeadlessRepo
		}
		err = EnsureCleanWorkingTree()
		if nil != err {
			return repo, err
		}
	}
	repo.LoadOrInit(directory)
	return repo, nil
}

func CommandInitAction(context *cli.Context) error {
	directory, _ := os.Getwd()
	repo := EnsureRepoIsAvailable(directory)
	return nil
}
