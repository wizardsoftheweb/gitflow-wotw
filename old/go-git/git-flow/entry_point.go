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
	VerbosityPattern, _ = regexp.Compile("^-v+$")
)

var (
	CliGlobalFlagVerbosity *cli.IntFlag
	CliFlagVersion         = cli.BoolFlag{
		Name:  "version, V",
		Usage: "Prints the version",
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
	CliGlobalFlagVerbosity = &cli.IntFlag{
		Name:  "verbose, v",
		Usage: "The more vs the more verbose the logs",
		Value: verbosity_level,
	}
	cli.VersionFlag = CliFlagVersion
	app := cli.NewApp()
	app.Name = "git-flow"
	app.Compiled = time.Now()
	app.Commands = []cli.Command{
		CommandInit,
	}
	app.Flags = []cli.Flag{
		CliGlobalFlagVerbosity,
	}
	app.Before = PopulateContext
	return app
}

func main() {
	// sanitized_args, verbosity_level := CheckVerbosity(os.Args)
	// app := BootstrapCli(verbosity_level)
	// err := app.Run(sanitized_args)
	// CheckError(err)
	directory, _ := os.Getwd()
	repo, _ := OpenRepoFromPath(directory)
	HasRemoteBranch(repo)
}
