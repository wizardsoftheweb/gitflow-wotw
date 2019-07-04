package main

import (
	"github.com/urfave/cli"
)

// func main() {

// }
//
var (
	TestCommand = cli.Command{
		Name:   "test",
		Flags:  []cli.Flag{},
		Action: ParsePrefix,
	}
)
