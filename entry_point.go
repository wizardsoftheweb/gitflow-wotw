package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

var (
	repo            Repository
	GITFLOW_VERSION = "0.0.0"
)

var (
	CliFlagVersion = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}
)

func BootstrapCli() *cli.App {
	app := cli.NewApp()
	app.Name = "git-flow"
	app.Version = GITFLOW_VERSION
	cli.VersionFlag = CliFlagVersion
	app.Compiled = time.Now()
	app.Commands = []cli.Command{
		CommandInit,
		CommandFeature,
		CommandHotfix,
		CommandRelease,
		CommandSupport,
		CommandVersion,
	}
	app.Flags = []cli.Flag{}
	return app
}

func main() {
	app := BootstrapCli()
	err := app.Run(os.Args)
	if nil != err {
		log.Fatal(err)
	}
}
