package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandHotfix = cli.Command{
		Name:   "hotfix",
		Flags:  []cli.Flag{},
		Action: CommandHotfixAction,
	}
)

func CommandHotfixAction(context *cli.Context) error {
	logrus.Trace("CommandHotfixAction")
	return nil
}
