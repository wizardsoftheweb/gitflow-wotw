package gitflow

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
	PackageCmd.AddCommand(InitCmd)
	InitCmd.Flags().BoolVarP(
		&Defaults,
		"defaults",
		"d",
		false,
		"Use defaults everywhere",
	)
	InitCmd.Flags().BoolVarP(
		&Force,
		"force",
		"f",
		false,
		"Force reinitialization",
	)
}
