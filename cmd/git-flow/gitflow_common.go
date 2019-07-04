package gitflow

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PassthroughThroughPrefixedBranchesWithErrorMessage(cmd *cobra.Command, args []string, remote bool) []string {
	branches := Repo.SpecificPrefixBranches(remote)
	if 0 == len(branches) {
		fmt.Println(
			fmt.Sprintf("There are no %s branches", Repo.HumanPrefix),
		)
		fmt.Println(
			fmt.Sprintf("The following command will set up a new %s branch:", Repo.HumanPrefix),
		)
		fmt.Println(
			fmt.Sprintf("\tgit flow %s start <name> [<base>]", Repo.HumanPrefix),
		)
	}
	return branches
}
