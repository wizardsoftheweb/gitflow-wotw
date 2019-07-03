package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	// https://stackoverflow.com/a/12093994
	// https://regex101.com/r/E2TCqU/3/tests
	GitReferenceRestrictionsPattern = regexp.MustCompile(`/\.|\.\.|\/\/+|[\000-\037\177 \\~^:?*[]+|(.(lock)?|/)$|^[^/]+$`)
	GitReferenceNotAt               = regexp.MustCompile(`^@$`)
	GitReferenceNoLeadingDots       = regexp.MustCompile(`^\.`)
	GitReferenceNoAtBracket         = regexp.MustCompile(`@\{`)
	GitReferenceNoSlashDot          = regexp.MustCompile(`/\.`)
	GitReferenceNoMultipleDot       = regexp.MustCompile(`\.\.+`)
	GitReferenceNoMultipleSlash     = regexp.MustCompile(`//+`)
	GitReferenceNoSpecialChars      = regexp.MustCompile(`[\000-\037\177 \~^:?*[]+`)
	GitReferenceNoDotLockSlashEnd   = regexp.MustCompile(`(\.(lock)?|/)$`)
	GitReferenceMustContainSlash    = regexp.MustCompile(`^[^/]+$`)
)

func ValidateRefName(ref_name string) error {
	if GitReferenceNotAt.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot only be '@': %s", ref_name))
	}
	if GitReferenceNoLeadingDots.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot lead with '.': %s", ref_name))
	}
	if GitReferenceNoAtBracket.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain '@{': %s", ref_name))
	}
	if GitReferenceNoSlashDot.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain '/.': %s", ref_name))
	}
	if GitReferenceNoMultipleDot.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain multiple consecutive '.': %s", ref_name))
	}
	if GitReferenceNoMultipleSlash.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain multiple consecutive '/': %s", ref_name))
	}
	if GitReferenceNoSpecialChars.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Contains an unallowed special character: %s", ref_name))
	}
	if GitReferenceNoDotLockSlashEnd.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot end in '.', '.lock', or '/': %s", ref_name))
	}
	if GitReferenceMustContainSlash.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Must contain at least one '/': %s", ref_name))
	}
	return nil
}

func ValidateBranchName(ref_name string) error {
	return ValidateRefName(fmt.Sprintf("refs/heads/%s", ref_name))
}

func PromptForBranchName(prompt_message string) string {
	prompt := promptui.Prompt{
		Label:    prompt_message,
		Validate: ValidateBranchName,
	}

	result, err := prompt.Run()
	if nil != err {
		fmt.Printf("Prompt failed %v\n", err)
		return PromptForBranchName(prompt_message)
	}
	return result
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
