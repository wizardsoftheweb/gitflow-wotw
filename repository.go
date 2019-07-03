package main

import (
	"errors"
	"log"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	GitBranchOutputPattern = regexp.MustCompile(`(?m)^\s*?(?P<branch>[^\s]+)\s*?$`)
)

var (
	ErrNotARepo = errors.New("Not a repo")
)

type Branch string

type Repository struct {
	dotDir        FileSystemObject
	configHandler ConfigFileHandler
	config        GitConfig
	localBranches []Branch

	Branches []Branch
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

func (repo *Repository) LoadLocalBranches() error {
	logrus.Trace("LoadLocalBranches")
	repo.localBranches = make([]Branch{})
	result := execute("git", "branch", "--no-color")
	for _, match := range GitBranchOutputPattern.FindAllStringSubmatch(result.stdout, -1) {
		result := map[string]string{}
		for index, name := range GitBranchOutputPattern.SubexpNames() {
			if 0 != index && "" != name {
				result[name] = string(match[index])
			}
		}
		repo.localBranches = append(branches, result["branch"])
	}
	return nil
}

func (repo *Repository) DoesBranchExistLocally(needle string) bool {
	for _, branch := range repo.localBranches {
		if branch == needle {
			return true
		}
	}
	return false
}

func (repo *Repository) HasBranchBeenConfigured(needle string) bool {
	branch_name := repo.config.Option(GIT_CONFIG_READ, "gitflow", "branch", branch)
	if "" != branch_name && repo.DoesBranchExistLocally(needle) {
		return true
	}
	return false
}
