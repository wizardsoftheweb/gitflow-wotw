package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var (
	VerbosityPattern = regexp.MustCompile("^-v+$")
)

var (
	GlobalCliFlagVerbosity *cli.IntFlag
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

func CheckVerbosity(args []string) ([]string, int) {
	sanitized_args := []string{}
	verbosity_flag := ""
	for _, arg := range os.Args {
		if VerbosityPattern.MatchString(arg) {
			verbosity_flag = fmt.Sprintf("%s%s", verbosity_flag, arg)
		} else {
			sanitized_args = append(sanitized_args, arg)
		}
	}
	return sanitized_args, strings.Count(verbosity_flag, "v")
}

func PopulateContext(context *cli.Context) error {
	BootstrapLogger(context.Int("verbose"))
	return nil
}

func BootstrapCli(verbosity_level int) *cli.App {
	GlobalCliFlagVerbosity = &cli.IntFlag{
		Name:  "verbose, v",
		Usage: "The more vs the more verbose the logs",
		Value: verbosity_level,
	}
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
	app.Flags = []cli.Flag{
		GlobalCliFlagVerbosity,
	}
	app.Before = PopulateContext
	return app
}

// func main() {
// 	sanitized_args, verbosity_level := CheckVerbosity(os.Args)
// 	app := BootstrapCli(verbosity_level)
// 	err := app.Run(sanitized_args)
// 	if nil != err {
// 		log.Fatal(err)
// 	}
// }
