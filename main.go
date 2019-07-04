package main

import (
	"fmt"
	"os"

	"github.com/wizardsoftheweb/gitflow-wotw/cmd"
)

func main() {
	if err := cmd.GitFlowCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
