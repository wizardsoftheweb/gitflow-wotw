package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnstagedChanges  = errors.New("There are unstaged changes in your working directory")
	ErrIndexUncommitted = errors.New("There are uncommitted changes in your index")
)

func IsRepoHeadless() bool {
	logrus.Trace("IsRepoHeadless")
	revparse := execute("git", "rev-parse", "--quiet", "--verify", "HEAD")
	return !revparse.Succeeded()
}

func IsWorkingTreeClean() bool {
	logrus.Trace("IsWorkingTreeClean")
	cleanliness := execute("git", "diff", "--no-ext-diff", "--ignore-submodules", "--quiet", "--exit-code")
	return cleanliness.Succeeded()
}

func IsIndexClean() bool {
	logrus.Trace("IsIndexClean")
	cleanliness := execute("git", "diff-index", "--cached", "--quiet", "--ignore-submodules", "HEAD", "--")
	return cleanliness.Succeeded()
}

func EnsureCleanWorkingTree() error {
	logrus.Trace("EnsureCleanWorkingTree")
	if !IsWorkingTreeClean() {
		return ErrUnstagedChanges
	} else {
		if !IsIndexClean() {
			return ErrIndexUncommitted
		}
	}
	return nil
}

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
