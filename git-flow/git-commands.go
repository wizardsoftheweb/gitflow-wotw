package main

import "gopkg.in/src-d/go-git.v4"

func GitInit(repo_path string) error {
	repo, err := git.PlainInit(repo_path)
	return err
}
