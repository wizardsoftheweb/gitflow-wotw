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

func InitProcedural(context *cli.Context) error {
	result := RevParseGitDir()
	if !result.Succeeded() {
		GitInit()
	} else {
		result = RevParseQuietVerifyHead()
		if !result.Succeeded() {
			IsWorkingTreeClean()
		} else {
			logrus.Fatal(ErrHeadlessRepo)
		}
	}
	if IsGitFlowInitialized() && !context.Bool("force") {
		logrus.Fatal(ErrAlreadyInitialized)
	}
	return nil
}

func CommandInitAction(context *cli.Context) error {
	logrus.Debug("CommandInitAction")
	return nil
}
