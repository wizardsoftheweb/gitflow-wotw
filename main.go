package main

import (
	"fmt"
	"os"

	"github.com/wizardsoftheweb/gitflow-wotw/cmd/gitflow"
)

func main() {
	if err := gitflow.GitFlowCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
