package main

import (
	"fmt"
	"os"

	gitflow "github.com/wizardsoftheweb/gitflow-wotw/cmd/git-flow"
)

func main() {
	if err := gitflow.PackageCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
