package main

import "github.com/urfave/cli"

var (
	CommandSupport = cli.Command{
		Name:   "support",
		Flags:  []cli.Flag{},
		Action: CommandSupportAction,
	}
)

func CommandSupportAction(context *cli.Context) error {
	return nil
}
