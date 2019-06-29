package main

import "os"

func main() {
	repo_path, _ := os.Getwd()
	GitFlowInit(repo_path)
}
