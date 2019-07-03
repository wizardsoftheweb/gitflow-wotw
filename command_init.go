package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	directory, _ := os.Getwd()
	dot_dir, _ := repo.discoverDotDir(FileSystemObject(directory))
	fmt.Println(repo.dotDir)
	config := &ConfigFileHandler{}
	config.configFile = FileSystemObject(filepath.Join(dot_dir.String(), "config"))
	config.loadConfig()
	config.parseConfig()
	return nil
}
