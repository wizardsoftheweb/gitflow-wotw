package main

import (
	"os"
)

var (
	repo Repository
)

func main() {
	cwd, _ := os.Getwd()
	repo.LoadOrInit(cwd)
	// path := FileSystemObject(cwd)
	// path.SearchDirectoryAbove("go-git2")
}
