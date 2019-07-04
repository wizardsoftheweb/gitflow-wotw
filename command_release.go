package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandRelease = cli.Command{
		Name:   "release",
		Flags:  []cli.Flag{},
		Action: CommandReleaseAction,
	}
)

func CommandReleaseAction(context *cli.Context) error {
	logrus.Debug("CommandReleaseAction")
	return nil
}
