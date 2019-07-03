package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

func GitInit(repo_path string) error {
	logrus.Debug("GitInit")
	_, err := git.PlainInit(repo_path, false)
	return err
}
