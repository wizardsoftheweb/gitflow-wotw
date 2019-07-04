package gitflow

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var FeatureCmd = &cobra.Command{
	Use:              "feature",
	TraverseChildren: true,
	PreRun:           PreRun,
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

func CommandFeatureListAction(cmd *cobra.Command, args []string) error {
	branches := PassthroughThroughPrefixedBranchesWithErrorMessage(cmd, args, false)
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
	return nil
}
