package gitflow

import (
	"github.com/spf13/cobra"
)

var PackageVersion = "0.0.0"

var PackageCmd = &cobra.Command{
	Use:     "git-flow",
	Version: PackageVersion,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var Verbosity int

func init() {
	PackageCmd.PersistentFlags().CountVarP(
		&Verbosity,
		"verbose",
		"v",
		"Increases application verbosity",
	)
}
