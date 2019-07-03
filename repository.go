package main

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	ErrNotARepo = errors.New("Not a repo")
)

type Repository struct {
	dotDir        FileSystemObject
	configHandler ConfigFileHandler
	config        GitConfig
}

func (repo Repository) discoverDotDir(root FileSystemObject) (FileSystemObject, error) {
	logrus.Trace("discoverDotDir")
	if root.isRoot() {
		return root, ErrNotARepo
	}
	result, err := root.DirectoryContains(".git")
	if nil != err {
		return root, err
	} else if result {
		return FileSystemObject(filepath.Join(root.String(), ".git")), nil
	}
	return repo.discoverDotDir(root.Parent())
}

func (repo *Repository) LoadConfig() error {
	var err error
	repo.configHandler.configFile = FileSystemObject(filepath.Join(repo.dotDir.String(), "config"))
	repo.configHandler.loadConfig()
	repo.config, err = repo.configHandler.parseConfig()
	if nil != err {
		log.Fatal(err)
	}
	return nil
}

func (repo *Repository) LoadOrInit(directory string) error {
	logrus.Trace("LoadOrInit")
	dot_dir, err := repo.discoverDotDir(FileSystemObject(directory))
	if nil != err {
		if ErrNotARepo == err {
			execute("git", "init")
		} else {
			return err
		}
	}
	repo.dotDir = dot_dir
	return nil
}
