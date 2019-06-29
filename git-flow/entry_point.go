package main

import "os"

func main() {
	BootstrapLogger()
	repo_path, _ := os.Getwd()
	GitFlowInit(repo_path)
}
