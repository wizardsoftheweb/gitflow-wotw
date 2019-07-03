package main

import "github.com/urfave/cli"

var (
	CommandRelease = cli.Command{
		Name:   "release",
		Flags:  []cli.Flag{},
		Action: CommandReleaseAction,
	}
)

func CommandReleaseAction(context *cli.Context) error {
	return nil
}
