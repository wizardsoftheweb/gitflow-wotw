package main

import (
	"regexp"

	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	// https://stackoverflow.com/a/12093994
	// https://regex101.com/r/E2TCqU/3/tests
	GitReferenceRestrictionsPattern = regexp.MustCompile(`^@$|^\.|@\{|\\|/\.|\.\.|\/\/+|[\000-\037\177 ~^:?*[]+|(.(lock)?|/)$|^[^/]+$`)
)

func ValidateRefName(ref_name string) bool {
	return !GitReferenceRestrictionsPattern.MatchString(ref_name)
}

func OpenRepoFromPath(repo_path string) (*git.Repository, error) {
	logrus.Debug("OpenRepoFromPath")
	repo, err := git.PlainOpenWithOptions(repo_path, &git.PlainOpenOptions{DetectDotGit: true})
	if nil != err {
		return nil, err
	}
	return repo, nil
}

func IsRepoHeadless(repo *git.Repository) bool {
	logrus.Debug("IsRepoHeadless")
	_, err := repo.ResolveRevision(plumbing.Revision(plumbing.HEAD))
	if plumbing.ErrReferenceNotFound == err {
		return true
	}
	return false
}

func GetSubmoduleNames(work_tree *git.Worktree) []string {
	logrus.Debug("GetSubmoduleNames")
	submodules, err := work_tree.Submodules()
	CheckError(err)
	names := make([]string, len(submodules))
	for index, submodule := range submodules {
		names[index] = submodule.Config().Path
	}
	return names
}

func AreThereUnstagedChanges(repo *git.Repository, ignore_submodules bool) bool {
	logrus.Debug("AreThereUnstagedChanges")
	work_tree, err := repo.Worktree()
	CheckError(err)
	changes, err := work_tree.Status()
	CheckError(err)
	files := make([]string, len(changes))
	index := 0
	for file := range changes {
		files[index] = file
		index++
	}
	if ignore_submodules {
		files, _ = RemoveStringElementFromStringSlice(files, ".gitmodules")
		for _, name := range GetSubmoduleNames(work_tree) {
			files, _ = RemoveStringElementFromStringSlice(files, name)
		}
	}
	return 0 != len(files)
}

func GetLocalBranchNames(repo_config *config.Config) []string {
	names := make([]string, len(repo_config.Branches))
	index := 0
	for branch_name, _ := range repo_config.Branches {
		names[index] = branch_name
		index++
	}
	return names
}

func IsRefNameValid(ref_name string) bool {
	return !GitReferenceRestrictionsPattern.Match(ref_name)
}
