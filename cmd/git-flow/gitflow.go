package gitflow

import (
	"github.com/spf13/cobra"
)

var PackageVersion = "0.0.0"

var PackageCmd = &cobra.Command{
	Use:              "git-flow",
	TraverseChildren: true,
	Version:          PackageVersion,
	PreRun: func(cmd *cobra.Command, args []string) {
		BootstrapLogger(VerbosityFlagValue)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var VerbosityFlagValue int

func init() {
	PackageCmd.PersistentFlags().CountVarP(
		&VerbosityFlagValue,
		"verbose",
		"v",
		"Increases application verbosity",
	)
}
