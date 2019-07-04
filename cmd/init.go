package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wat")
	},
}

var Defaults bool
var Force bool

func init() {
	GitFlowCmd.AddCommand(InitCmd)
	InitCmd.Flags().BoolVarP(
		&Defaults,
		"defaults",
		"d",
		false,
		"Use defaults everywhere",
	)
}
