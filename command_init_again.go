package main

import (
	"github.com/urfave/cli"
)

// func main() {

// }
//
var (
	TestCommand = cli.Command{
		Name:   "test",
		Flags:  []cli.Flag{},
		Action: ParsePrefix,
	}
)

func ParsePrefix(context *cli.Context, prefixKey string, defaultValue string) error {
	if !context.Bool("force") {
		return nil
	} else {
		var newValue string
		if context.Bool("default") {
			newValue = defaultValue
		} else {
			value := GitConfig.GetWithDefault(prefixKey, defaultValue)
			newValue = PromptForInput(PREFIX_NAME_VALIDATION, "test", value)
		}
		GitConfig.Write(prefixKey, newValue)
	}
	return nil
}
