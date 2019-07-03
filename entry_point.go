package main

import "os"

var (
	repo Repository
)

func main() {
	cwd, _ := os.Getwd()
	repo.LoadOrInit(cwd)
}
