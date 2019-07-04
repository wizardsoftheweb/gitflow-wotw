package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandSupport = cli.Command{
		Name:   "support",
		Flags:  []cli.Flag{},
		Action: CommandSupportAction,
	}
)

func CommandSupportAction(context *cli.Context) error {
	logrus.Debug("CommandSupportAction")
	return nil
}
