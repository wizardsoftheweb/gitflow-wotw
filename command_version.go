package main

import "github.com/urfave/cli"

var (
	CommandVersion = cli.Command{
		Name:   "version",
		Flags:  []cli.Flag{},
		Action: CommandVersionAction,
	}
)

func CommandVersionAction(context *cli.Context) error {
	return nil
}
