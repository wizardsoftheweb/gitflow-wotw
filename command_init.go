package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

var (
	CommandInit = cli.Command{
		Name:   "init",
		Flags:  []cli.Flag{},
		Action: CommandInitAction,
	}
)

func CommandInitAction(context *cli.Context) error {
	repo := &Repository{}
	// directory, _ := os.Getwd()
	// dot_dir, _ := repo.discoverDotDir(FileSystemObject(directory))
	fmt.Println(repo.dotDir)
	config := &ConfigEnvironmentHandler{}
	// config.configFile = FileSystemObject(filepath.Join(dot_dir.String(), "config"))
	config.loadConfig()
	loadedConfig, err := config.parseConfig()
	if nil != err {
		log.Fatal(err)
	}
	for _, line := range config.dumpConfig(loadedConfig) {
		fmt.Println(line)
	}
	return nil
}
