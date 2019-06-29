package main

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	format "gopkg.in/src-d/go-git.v4/plumbing/format/config"
)

func OpenRepoFromPath(repo_path string) (*git.Repository, error) {
	repo, err := git.PlainOpenWithOptions(repo_path, &git.PlainOpenOptions{DetectDotGit: true})
	if nil != err {
		return nil, err
	}
	return repo, nil
}

func IsRepoHeadless(repo *git.Repository) bool {
	_, err := repo.ResolveRevision(plumbing.Revision(plumbing.HEAD))
	if plumbing.ErrReferenceNotFound == err {
		return true
	}
	return false
}

func GetSubmoduleNames(work_tree *git.Worktree) []string {
	submodules, err := work_tree.Submodules()
	CheckError(err)
	names := make([]string, len(submodules))
	for index, submodule := range submodules {
		names[index] = submodule.Config().Path
	}
	return names
}

func AreThereUnstagedChanges(repo *git.Repository, ignore_submodules bool) bool {
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

func GetConfigOptions(options format.Options) {
	for _, option := range options {
		fmt.Println(option.Key, option.Value)
	}
}

func GetConfigValue(repo *git.Repository) interface{} {
	config, err := repo.Config()
	CheckError(err)
	for _, section := range config.Raw.Sections {
		fmt.Println(section.Name, section)
		GetConfigOptions(section.Options)
		for _, subsection := range section.Subsections {
			fmt.Println(subsection.Name)
			GetConfigOptions(subsection.Options)
		}
	}
	fmt.Println("rad")
	return nil
}
