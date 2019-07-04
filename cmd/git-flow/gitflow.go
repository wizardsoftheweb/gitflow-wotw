package gitflow

import (
	"github.com/spf13/cobra"
)

var GITFLOW_VERSION = "0.0.0"

var GitFlowCmd = &cobra.Command{
	Use:     "gitflow",
	Version: GITFLOW_VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var Verbosity int

func init() {
	GitFlowCmd.PersistentFlags().CountVarP(
		&Verbosity,
		"verbose",
		"v",
		"Increases application verbosity",
	)
}
