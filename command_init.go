package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandInit = cli.Command{
		Name:   "release",
		Flags:  []cli.Flag{},
		Action: CommandInitAction,
	}
)

func CommandInitAction(context *cli.Context) error {
	logrus.Debug("CommandInitAction")
	return nil
}
