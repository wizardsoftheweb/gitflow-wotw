package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandFeature = cli.Command{
		Name:   "feature",
		Flags:  []cli.Flag{},
		Action: CommandFeatureAction,
	}
)

func CommandFeatureAction(context *cli.Context) error {
	logrus.Trace("CommandFeatureAction")
	return nil
}
