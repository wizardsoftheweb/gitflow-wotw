package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func CheckError(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

func IsCwdInRepo() bool {
	current_path, _ := os.Getwd()
	_, err := OpenRepoFromPath(current_path)
	if git.ErrRepositoryNotExists == err {
		return false
	}
	return true
}

func GitFlowInit() {
	current_path, _ := os.Getwd()
	repo, err := OpenRepoFromPath(current_path)
	if git.ErrRepositoryNotExists == err {
		fmt.Println("not a repo")
	} else {
		IsRepoHeadless(repo)
		AreThereUnstagedChanges(repo, true)
	}
}
