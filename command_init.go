package main

import "github.com/urfave/cli"

var (
	CommandInit = cli.Command{
		Name:   "init",
		Flags:  []cli.Flag{},
		Action: CommandInitAction,
	}
)

func CommandInitAction(context *cli.Context) error {
	return nil
}
