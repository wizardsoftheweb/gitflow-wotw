package main

import (
	"fmt"

	"github.com/urfave/cli"
)

const (
	defaultGitFlowRootPackage = "gitflow"
	defaultGitFlowRootVersion = "0.4.2-pre"
	defaultGitFlowRootUrl     = "https://github.com/nvie/gitflow"
)

var (
	CommandVersion = cli.Command{
		Name:   "version",
		Flags:  []cli.Flag{},
		Action: CommandVersionAction,
	}
)

func CommandVersionAction(context *cli.Context) error {
	fmt.Fprintf(
		context.App.Writer,
		"Version %s of gitflow-wotw was based on the following work:\n\tPackage: %s\n\tVersion: %s\n\tUrl: %s",
		context.App.Version,
		defaultGitFlowRootPackage,
		defaultGitFlowRootVersion,
		defaultGitFlowRootUrl,
	)
	return nil
}
