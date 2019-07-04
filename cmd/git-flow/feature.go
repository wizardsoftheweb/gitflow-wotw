package gitflow

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var FeatureCmd = &cobra.Command{
	Use:    "feature",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wat")
	},
}

func init() {
	PackageCmd.AddCommand(FeatureCmd)
}

func PreRun(cmd *cobra.Command, args []string) {
	Repo.Prefix = GitConfig.GetWithDefault(FeaturePrefixKey, DefaultPrefixFeature.Value)
	Repo.HumanPrefix = strings.TrimSuffix(Repo.Prefix, "/")
}
