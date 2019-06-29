package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
)

func BootstrapCli() {
	app := cli.NewApp()
	app.Name = "git-flow"
	app.Compiled = time.Now()
}

func main() {
	BootstrapLogger()
	repo_path, _ := os.Getwd()
	GitFlowInit(repo_path)
}
