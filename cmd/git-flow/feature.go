package gitflow

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	PackageCmd.AddCommand(FeatureCmd)
	FeatureCmd.AddCommand(FeatureListAction)
}

var FeatureCmd = &cobra.Command{
	Use:              "feature",
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Repo.Prefix = GitConfig.GetWithDefault(FeaturePrefixKey, DefaultPrefixFeature.Value)
		Repo.HumanPrefix = strings.TrimSuffix(Repo.Prefix, "/")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wat")
	},
}

var FeatureListAction = &cobra.Command{
	Use:              "list",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		branches := PassthroughThroughPrefixedBranchesWithErrorMessage(false)
		logrus.Trace("message")
		var width int
		if 0 < VerbosityFlagValue {
			for _, branch := range branches {
				width = MaxInt(width, len(branch))
			}
			width += 3
		}
		for _, branch := range branches {
			branchFullname := fmt.Sprintf("%s%s", Repo.Prefix, branch)
			if Repo.CurrentBranch() == branchFullname {
				fmt.Printf("* ")
			} else {
				fmt.Printf("  ")
			}
			if 0 < VerbosityFlagValue {
				devBranch := GitConfig.Get(DevBranchKey)
				base := MergeBase(branchFullname, devBranch)
				devSha := RevParseArgs(devBranch)
				branchSha := RevParseArgs(branchFullname)
				fmt.Printf(fmt.Sprintf("%%-%ds", width), branch)
				switch {
				case branchSha == devSha:
					fmt.Printf("(sha is identical to dev)")
				case base == branchSha:
					fmt.Printf("(may ff to dev)")
				case base == devSha:
					fmt.Printf("(identical to latest dev")
				default:
					fmt.Printf("may be rebased")
				}

			} else {
				fmt.Printf("%s", branch)
			}
			fmt.Println()
		}
	},
}
