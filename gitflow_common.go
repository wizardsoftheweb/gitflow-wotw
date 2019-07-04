package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func PassthroughThroughPrefixedBranchesWithErrorMessage(context *cli.Context, remote bool) []string {
	branches := Repo.SpecificPrefixBranches(remote)
	if 0 == len(branches) {
		fmt.Fprintln(
			context.App.Writer,
			fmt.Sprintf("There are no %s branches", Repo.HumanPrefix),
		)
		fmt.Fprintln(
			context.App.Writer,
			fmt.Sprintf("The following command will set up a new %s branch:", Repo.HumanPrefix),
		)
		fmt.Fprintln(
			context.App.Writer,
			fmt.Sprintf("\tgit flow %s start <name> [<base>]", Repo.HumanPrefix),
		)
	}
	return branches
}
