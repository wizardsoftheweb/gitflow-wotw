package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	ErrUnableToGitInit    = errors.New("Unable to complete git init in the current working directory")
	ErrHeadlessRepo       = errors.New("Unable to initialize in a bare repo")
	ErrAlreadyInitialized = errors.New("The repo is already initialized; try again with -f")
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
	logrus.Trace("EnsureRepoIsAvailable")
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
		if !headless {
			err = EnsureCleanWorkingTree()
			if nil != err {
				return repo, err
			}
		}

	}
	repo.LoadOrInit(directory)
	return repo, nil
}

func CheckInitialization(context *cli.Context, repo Repository) error {
	logrus.Trace("CheckInitialization")
	if IsGitflowInitialized(repo.config) && context.Bool("force") {
		return ErrAlreadyInitialized
	}
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Trace("CommandInitAction")
	directory, _ := os.Getwd()
	repo, err := EnsureRepoIsAvailable(directory)
	if nil != err {
		log.Fatal(err)
	}
	if context.Bool("default") {
		logrus.Info("Using default branches")
	}
	fmt.Println(repo)
	return nil
}
